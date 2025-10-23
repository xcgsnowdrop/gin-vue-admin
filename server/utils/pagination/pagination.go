package pagination

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PageResult 分页结果
type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
}

// PageParam 分页参数
type PageParam struct {
	PageNum  int64 // 页码，从1开始
	PageSize int64 // 每页数量
}

// FindOptions 查询选项
type FindOptions struct {
	Projection interface{} // 投影字段
	Sort       interface{} // 排序条件
}

// FindWithPage 通用MongoDB分页查询
//   - collection: mongo集合
//   - filter: 查询条件
//   - page: 分页参数
//   - sort: 排序条件（如 bson.D{{"createTime", -1}}）
//   - resultMaker: func() T 用于创建新对象（如 func() *DbGuild { return &DbGuild{} }）
//   - projection: 可选，指定返回字段（如 bson.M{"_id": 0, "name": 1}）
//
// 返回 PageResult[T] 和 error
//
// resultMaker 参数说明：
//   - resultMaker 是一个工厂函数，用于动态创建一个 T 类型的实例（通常是指针），以便 MongoDB 的 Decode 方法能正确填充数据。
//   - 由于 Go 泛型的限制，无法在泛型函数内部直接 new 出正确类型的指针（如 *DbGuild），也无法自动推断 T 的具体类型。
//   - MongoDB 的 Decode 方法需要一个非 nil 的指针类型目标对象，否则会解码失败或 panic。
//   - 只有调用者最清楚 T 的具体类型和如何正确初始化，因此必须由调用者传入 resultMaker。
//
// 为什么不能省略 resultMaker：
//   - Go 泛型不支持像 C++/Java 那样的"自动 new T()"，new(T) 得到的是 **DbGuild 这类类型，类型不匹配。
//   - 直接 var item T 得到的是 nil 指针，不能用于 Decode。
//   - 省略 resultMaker 会导致类型不安全或运行时错误。
//
// 结论：
//   - resultMaker 是为了解决 Go 泛型和 MongoDB Decode 结合时的类型安全和实例化问题，不能省略，是实现通用分页的最佳实践。
func FindWithPage[T any](
	ctx context.Context,
	collection *mongo.Collection,
	filter interface{},
	page PageParam,
	sort interface{},
	resultMaker func() T,
	projection ...interface{},
) (PageResult[T], error) {
	var res PageResult[T]

	// 统计总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return res, err
	}
	res.Total = total

	// 查询分页数据
	opts := options.Find().
		SetSort(sort).
		SetSkip((page.PageNum - 1) * page.PageSize).
		SetLimit(page.PageSize)

	// 如果提供了 projection，则设置投影
	if len(projection) > 0 && projection[0] != nil {
		opts.SetProjection(projection[0])
	}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return res, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		item := resultMaker()
		if err := cursor.Decode(item); err != nil {
			return res, err
		}
		res.List = append(res.List, item)
	}
	return res, nil
}

// FindWithPageOptions 使用选项结构体的分页查询（推荐使用）
func FindWithPageOptions[T any](
	ctx context.Context,
	collection *mongo.Collection,
	filter interface{},
	page PageParam,
	resultMaker func() T,
	findOptions *FindOptions,
) (PageResult[T], error) {
	var res PageResult[T]

	// 统计总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return res, err
	}
	res.Total = total

	// 查询分页数据
	opts := options.Find().
		SetSkip((page.PageNum - 1) * page.PageSize).
		SetLimit(page.PageSize)

	// 设置排序
	if findOptions != nil && findOptions.Sort != nil {
		opts.SetSort(findOptions.Sort)
	}

	// 设置投影
	if findOptions != nil && findOptions.Projection != nil {
		opts.SetProjection(findOptions.Projection)
	}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return res, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		item := resultMaker()
		if err := cursor.Decode(item); err != nil {
			return res, err
		}
		res.List = append(res.List, item)
	}
	return res, nil
}
