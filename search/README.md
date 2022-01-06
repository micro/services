Indexing and full text search

# Search Service

Store and search JSON documents. The Search API provides full indexing and text search.

Powered by [OpenSearch](https://opensearch.org/).

Search for a given word or phrase in a particular field of a document. Combine multiple with either `AND` or `OR` boolean operators to create complex queries.

## Usage
Documents are inserted using the `/search/index` endpoint. Document fields are automatically indexed with no need to define which fields to index ahead of time. Documents are logically grouped in to `indexes` so you may have an index for customers and one for products. Once documents are inserted you are ready to search, simple as that.

## Search query language

The search API supports a simple query language to let you get to your data quickly without having to learn a complicated language. 

The most basic query looks like this

```sql
key == 'value'
```

where you specify a key and a value to find. For example you might want to look for every customer with first name of John

```sql
first_name == 'John'
```

String values support single or double quotes. 

Values can also be numbers 

```sql
age == 37
```

or booleans

```sql
verified == true
```

You can search on fields that are nested in the document using dot (`.`) as a separator 

```sql
address.city == 'London'
```

In addition to equality `==` the API support greater than or equals `>=` and less than or equals `<=` operators

```sql
age >= 37
age <= 37
```

Simple queries can be combined with logical `and` 

```sql
first_name == "John" AND age <= 37
```

or logical `or`
```sql
first_name == "John" OR first_name == "Jane"
```

If combining `and` and `or` operations you will need to use parentheses to explicitly define precedence

```sql
(first_name == "John" OR first_name == "Jane") AND age <= 37 
```
