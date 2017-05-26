# csviewer

csviewer is command line csv viewer.


# Install

```
go get github.com/Konboi/csviewer
```
### From File

using `p` or `path` option

```
csviewer -p _example/example.csv
```

### From stdin

```
cat _example/example.csv | csviewer
```
## Display Option

### Default

```
 $ csviewer -p _example/example.csv
+----+------+----------------+-----------+--------+
| ID | NAME |      MAIL      |   PHONE   | ADRESS |
+----+------+----------------+-----------+--------+
|  1 | a    | aaaa@hoge.fuga |    123456 | 111111 |
|  2 | b    | bbb@hoge.fuga  |     12345 |        |
|  3 | c    |                |           |  22222 |
|  5 | d    | ddd@fuga.hgoe  | 123456789 |        |
+----+------+----------------+-----------+--------+
```

### Columns option

set display columns.

```
 $ csviewer -p _example/example.csv -c id,mail,name
+----+----------------+------+
| ID |      MAIL      | NAME |
+----+----------------+------+
|  1 | aaaa@hoge.fuga | a    |
|  2 | bbb@hoge.fuga  | b    |
|  3 |                | c    |
|  5 | ddd@fuga.hgoe  | d    |
+----+----------------+------+
```

### Limit Option

set display rows num.

```
 $ csviewer -p _example/example.csv -l 2
+----+------+----------------+--------+--------+
| ID | NAME |      MAIL      | PHONE  | ADRESS |
+----+------+----------------+--------+--------+
|  1 | a    | aaaa@hoge.fuga | 123456 | 111111 |
|  2 | b    | bbb@hoge.fuga  |  12345 |        |
+----+------+----------------+--------+--------+
```

### Filter Option

set display conition.

```
 $ ./csviewer -p _example/example.csv -f "id > 2"
+----+------+---------------+-----------+--------+
| ID | NAME |     MAIL      |   PHONE   | ADRESS |
+----+------+---------------+-----------+--------+
|  3 | c    |               |           |  22222 |
|  5 | d    | ddd@fuga.hgoe | 123456789 |        |
+----+------+---------------+-----------+--------+
```

```
 $ csviewer  -p _example/example.csv -f 'phone > 12345'
+----+------+----------------+-----------+--------+
| ID | NAME |      MAIL      |   PHONE   | ADRESS |
+----+------+----------------+-----------+--------+
|  1 | a    | aaaa@hoge.fuga |    123456 | 111111 |
|  5 | d    | ddd@fuga.hgoe  | 123456789 |        |
```

# Useage


```
$ csviewer --help
Usage of ./csviewer:
  -c string
        print specify columns
  -columns string
        print specify columns
  -f value
        filter
  -filter value
        filter
  -l int
        set max display rows num
  -limit int
        set max display rows  num
  -p string
        set csv file path
  -path string
        set csv file path
```

# TODO

- [ ] Order option
- [ ] Set multi filters in one column