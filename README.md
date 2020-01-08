# gourm

### A simple golang orm

##### 3.add update, which function support column and value mode

##### 2.now only postgresSQL

##### 1.first commit

## simple usage

### struct simple

### basic functions

``` go
type User struct {
    gourm.Model `table:"xxx" primary_key:"id"` 
        ID        int `col:"id"` 
        Name      string `col:"name"` 
        Age       int `col:"age"` 
        LoginTime string `col:"string"` 
}

```

``` go
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

### 3 - updates:

``` go
db.Which(&struct, col string, val interface{})
db.Update(&struct, col string, val interface{})

```
