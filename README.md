# gourm
### A simple golang orm

##### first commit
##### now only postgresSQL

### simple usage

```go

db := gourm.New("postgres", dbconfig, ifping)

db.Insert(&struct)

db.Update(&struct)
    db.Select("name", "age").Where("age < ?", 20).Update(&struct)

db.Which(&struct)
    db.Which(&struct, "name", "xxx")

    db.Where("name = ?", "xxx").Find(&[]struct)
    db.Where("id <> ?", 0).Order("name desc").Offset(10).Limit(5).Find(&[]struct)

db.Delete(&struct)
    db.Delete(&struct, "id", 9)

    ```
