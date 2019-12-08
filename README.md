# psql-splitter

Split a (postgres) sql file into multiple files by a chunk of statements.

## Installation

`go get github.com/nrnrk/psql-splitter`

## Usage

`psql-splitter split {sql file} [-n {split by}]`

For example, if you split `sample.sql` like

```
SELECT * FROM users;
SELECT * FROM others;
INSERT INTO users
    VALUES (12111, 'Mike', TRUE);
```

by 2 statements, execute the following command.

`psql-splitter sample.sql -n 2`

The 2 following files are generated.

sample-aa.sql
```
SELECT * FROM users;
SELECT * FROM others;
```

sample-ab.sql
```

INSERT INTO users
    VALUES (12111, 'Mike', TRUE);
```

You can confirm this split by `diff` command like

```
cat sample-* |diff sample.sql -
```

## License

This software includes the work that is distributed in the [Apache License 2.0](http://www.apache.org/licenses/LICENSE-2.0).
