package mongodb

import (
	"context"
	"math"
	"reflect"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DB(mongo *mongo.Database) *mongo.Database {
	return mongo
}

type Repo struct {
	Collection *mongo.Collection
	Mux        sync.Mutex
}

func (r *Repo) Create(i interface{}) error {
	rt := reflect.TypeOf(i)
	//print(rt.Kind())
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		return r.createMany(i)
	default:
		return r.createOne(i)
	}
}

func (r *Repo) createOne(i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if m, ok := i.(ModelInterface); ok {
		m.Stamp()
		m.SetID(primitive.NewObjectID())
	}
	if _, err := r.Collection.InsertOne(ctx, i); err != nil {
		return err
	}
	return nil
}

func (r *Repo) UpdateOneByPrimitiveM(id primitive.ObjectID, u primitive.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	r.Mux.Lock()
	_, err := r.Collection.UpdateOne(ctx,
		primitive.D{
			primitive.E{
				Key:   "_id",
				Value: id,
			},
		}, u)
	r.Mux.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) createMany(i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	v := reflect.ValueOf(i)
	is := []interface{}{}
	for i := 0; i < v.Len(); i++ {
		if m, ok := v.Index(i).Interface().(ModelInterface); ok {
			m.Stamp()
			m.SetID(primitive.NewObjectID())
			is = append(is, m)
		}
	}
	if _, err := r.Collection.InsertMany(ctx, is); err != nil {
		return err
	}
	return nil
}

// Update update
func (r *Repo) Update(i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var id primitive.ObjectID
	if m, ok := i.(ModelInterface); ok {
		m.UpdateStamp()
		id = m.GetID()
	}
	r.Mux.Lock()
	err := r.Collection.FindOneAndReplace(ctx, primitive.M{"_id": id}, i).Err()
	r.Mux.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) UpdateUpsert(filter primitive.M, update primitive.M, i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if m, ok := i.(ModelInterface); ok {
		m.UpdateStamp()
	}
	opts := options.Update().SetUpsert(true)
	r.Mux.Lock()
	_, err := r.Collection.UpdateOne(ctx, filter, update, opts)
	r.Mux.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) Delete(i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var id primitive.ObjectID
	if m, ok := i.(ModelInterface); ok {
		m.UpdateDeletedStamp()
		id = m.GetID()
	}
	r.Mux.Lock()
	_, err := r.Collection.UpdateOne(ctx,
		primitive.D{
			primitive.E{
				Key:   "_id",
				Value: id,
			},
		}, primitive.M{
			"$set": i,
		})
	r.Mux.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) HardDelete(i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var id primitive.ObjectID
	if m, ok := i.(ModelInterface); ok {
		id = m.GetID()
	}
	r.Mux.Lock()
	_, err := r.Collection.DeleteOne(ctx,
		primitive.D{
			primitive.E{
				Key:   "_id",
				Value: id,
			},
		})
	r.Mux.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) FindOneByPrimitiveM(m primitive.M, i interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := r.Collection.FindOne(ctx, m).Decode(i)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) FindOneByID(id string, i interface{}) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = r.FindOneByPrimitiveM(primitive.M{
		"_id": oid,
	}, i)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) FindOneByObjectID(oid *primitive.ObjectID, i interface{}) error {
	err := r.FindOneByPrimitiveM(primitive.M{
		"_id": oid,
	}, i)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) FindAllByPrimitiveM(m primitive.M, result interface{}, opts ...*options.FindOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := r.Collection.Find(ctx, m, opts...)
	if err != nil {
		return err
	}
	defer func() { _ = cur.Close(ctx) }()
	if err := r.BindModels(ctx, cur, result); err != nil {
		return err
	}
	return nil
}

func (r *Repo) BindModels(ctx context.Context, cur *mongo.Cursor, result interface{}) error {
	resultv := reflect.ValueOf(result)
	slicev := resultv.Elem()
	if slicev.Kind() == reflect.Interface {
		slicev = slicev.Elem()
	}
	slicev = slicev.Slice(0, slicev.Cap())
	elemt := slicev.Type().Elem()
	i := 0

	for {
		elemp := reflect.New(elemt)
		if !cur.Next(ctx) {
			break
		}
		err := cur.Decode(elemp.Interface())
		if err != nil {
			return err
		}
		slicev = reflect.Append(slicev, elemp.Elem())
		i++
	}
	resultv.Elem().Set(slicev.Slice(0, i))
	return nil
}

func (r *Repo) Aggregate(pipeline []primitive.M, i interface{}, form ...*PageQuery) (*Page, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var size int64
	var page int64
	if len(form) <= 0 {
		form := &PageQuery{}
		size = form.GetSize()
		page = form.GetPage()
	} else {
		size = form[0].GetSize()
		page = form[0].GetPage()
	}
	//pipelineCount := pipeline
	//=====================Count=========================
	//###################Channels###################
	ch := make(chan []*CountDocument) // สร้างท่อ Channel เอาไว้ส่งข้อมูล
	go r.countDocumentByAggregate(pipeline, ch, ctx)

	//#############################################
	/*pipelineCount = append(pipelineCount, primitive.M{
		"$count": "count",
	})
	cursorCount, err := r.Collection.Aggregate(
		ctx,
		pipelineCount,
	)
	if err != nil {
		return nil, err
	}
	countDocument := []*CountDocument{}
	if err = cursorCount.All(ctx, &countDocument); err != nil {
		return nil, err
	}*/
	//log.Println("countDocument:",countDocument[0].Count)
	//===================================================
	//###################################################
	pipeline = append(pipeline, primitive.M{
		"$skip": int64(size * (page - 1)),
	})
	pipeline = append(pipeline, primitive.M{
		"$limit": int64(size),
	})
	cursor, err := r.Collection.Aggregate(
		ctx,
		pipeline,
	)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, i); err != nil {
		return nil, err
	}
	//###################################################
	countDocument := <-ch // ค่าจากท่อ Channel จะออกตรงนี้
	var p Page
	if len(countDocument) > 0 {
		count := countDocument[0].Count
		if count > 0 {
			p.PageInformation.Page = page
			p.PageInformation.Size = size
			p.PageInformation.TotalNumberOfEntities = count
			p.PageInformation.TotalNumberOfPages = int64(math.Ceil(float64(count) / float64(size)))
			p.Entities = i
			return &p, nil
		}

	}
	emptyArray := make([]interface{}, 0)
	p.Entities = emptyArray
	return &p, nil
}

func (r *Repo) countDocumentByAggregate(pipeline []primitive.M, data chan []*CountDocument, ctx context.Context) error {
	pipeline = append(pipeline, primitive.M{
		"$count": "count",
	})
	cursorCount, err := r.Collection.Aggregate(
		ctx,
		pipeline,
	)
	if err != nil {
		return err
	}
	countDocument := []*CountDocument{}
	if err = cursorCount.All(ctx, &countDocument); err != nil {
		return err
	}
	data <- countDocument // นำค่าเข้าท่อ Channel
	return nil
}
