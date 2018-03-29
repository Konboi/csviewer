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

set display condition.

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

### Multiple Filter Option

#### And

```
$ ./csviewer -p _example/example.csv -f "id > 2 && id <= 10"
+----+------+-----------------+-----------+--------+
| ID | NAME |      MAIL       |   PHONE   | ADRESS |
+----+------+-----------------+-----------+--------+
|  3 | c    |                 |           |  22222 |
|  5 | d    | ddd@fuga.hgoe   | 123456789 |        |
| 10 | e    | eeeee@fuga.hgoe |    654321 |        |
+----+------+-----------------+-----------+--------+
```

```
$ ./csviewer -p _example/example.csv -f "id > 2" -f "id <= 10"
+----+------+-----------------+-----------+--------+
| ID | NAME |      MAIL       |   PHONE   | ADRESS |
+----+------+-----------------+-----------+--------+
|  3 | c    |                 |           |  22222 |
|  5 | d    | ddd@fuga.hgoe   | 123456789 |        |
| 10 | e    | eeeee@fuga.hgoe |    654321 |        |
+----+------+-----------------+-----------+--------+
```

#### Or


```
$ ./csviewer -p _example/example.csv -f "name == 'c'" -f "name == 'd'" -or
+----+------+---------------+-----------+--------+
| ID | NAME |     MAIL      |   PHONE   | ADRESS |
+----+------+---------------+-----------+--------+
|  3 | c    |               |           |  22222 |
|  5 | d    | ddd@fuga.hgoe | 123456789 |        |
+----+------+---------------+-----------+--------+
```

```
$ ./csviewer -p _example/example.csv -f "name == 'c' || name == 'd'"
+----+------+---------------+-----------+--------+
| ID | NAME |     MAIL      |   PHONE   | ADRESS |
+----+------+---------------+-----------+--------+
|  3 | c    |               |           |  22222 |
|  5 | d    | ddd@fuga.hgoe | 123456789 |        |
+----+------+---------------+-----------+--------+
```

### Sort Option

```
 $ csviewer -p _example/example.csv -s 'phone asc'
+-----+-------+-----------------+-----------+--------+
| ID  | NAME  |      MAIL       |   PHONE   | ADRESS |
+-----+-------+-----------------+-----------+--------+
|   3 | c     |                 |           |  22222 |
|   2 | b     | bbb@hoge.fuga   |     12345 |        |
|   1 | a     | aaaa@hoge.fuga  |    123456 | 111111 |
|  10 | e     | eeeee@fuga.hgoe |    654321 |        |
|   5 | d     | ddd@fuga.hgoe   | 123456789 |        |
| 222 | asdfg | asdfg@fuga.hgoe | 987654321 |        |
+-----+-------+-----------------+-----------+--------+
```

```
 $ csviewer -p _example/example.csv -s 'mail desc'
+-----+-------+-----------------+-----------+--------+
| ID  | NAME  |      MAIL       |   PHONE   | ADRESS |
+-----+-------+-----------------+-----------+--------+
|  10 | e     | eeeee@fuga.hgoe |    654321 |        |
|   5 | d     | ddd@fuga.hgoe   | 123456789 |        |
|   2 | b     | bbb@hoge.fuga   |     12345 |        |
| 222 | asdfg | asdfg@fuga.hgoe | 987654321 |        |
|   1 | a     | aaaa@hoge.fuga  |    123456 | 111111 |
|   3 | c     |                 |           |  22222 |
+-----+-------+-----------------+-----------+--------+
```

# Usage


```
 $ csviewer --help
Usage of csviewer:
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
  -s string
        sort by set value
        ex) id desc/ hoge_id asc
  -sort string
        sort by set value
        ex) id desc/ hoge_id asc
```

# TODO

- [x] Order option
- [x] Set multi filters in one column
