#### Why use an ORM with type checking?

With gorm you do queries like `db.Where("name = ?", name).Find(&user)` and the query must be unit tested as it could have 
spelling mistakes/copy-paste errors.

The entgo equivalent is `db.User.Query().Where(user.NameEQ(name)).First(context.TODO())` and we get
field type checking.
