# Language

## Repository
This project is hosted on GitHub
[here](https://github.com/Zac-Garby/language).

## Intro
Recently, I've been working on a programming language, in Go. Here's a basic
hello world program in it:

```go
print("Hello, world");
```

In it, I tried to add as many weird/obscure features as I could, which I'll
demonstrate later.

## Get it
It's quite easy to install:

```shell
$ git clone https://github.com/Zac-Garby/language.git
$ cd language
$ ./build/main
```

This will run the REPL. It works just like any other REPL - type an expression,
get the value. You can also execute a file:

```shell
$ ./build/main <filename>
```

Now have a look at the examples below to see some of the things you can do!

## Loops
Firstly, for and while loops return lists, containing their values at each
iteration:

```go
array := 1..10;

result := for (i | array) {
  array[i] ** 2;
};

print(result);
```
```shell
$ ./main

[2, 4, 6, 8, 10, 12, 14, 16, 18, 20]
```

`1..10` creates an array containing all the integers from 1 to 10 (inclusive.)
You can also do `1..<10` to create an array with all the integers from 1 (inclusive)
to 10 (not inclusive), so it's equivalent to `1..9`.

The for loop uses the syntax `for (counter | set) { body }`. The for loop is an
expression, which means it returns a value. The value is an array, containing
the value of *body* at each iteration. The value of *body* can be implicitly or
explicitly returned - you can either use a return statement or just have the value
as the result of the last expression in the body.

The same goes for while loops: `while (condition) { body }`.

You can of course use loop control statements like in any other language: `next`,
and `break`, and they do exactly what you'd expect.

## Object system
I also designed my own object system, similar to JavaScript's prototype based
"classes". In my language, they're called *models*.

To understand models, you first need to understand hashes. A hash is like a
JavaScript object: it has a mapping of string keys to values.

You can create a hash like this:

```go
a := {
  x: 2,
  y: 3,
  z: {
    hello: 6,
    world: \(x, y) = x + y
  }
};
```

As you can see, they can be n-dimensional. You can access them by string value
or identifier:

```go
print("a.x                =", a.x);
print("a['y']             =", a["y"]);
print("a['z'].world(2, 3) =", a["z"].world(2, 3));
```
```shell
$ ./main

a.x                = 2
a['y']             = 3
a['z'].world(2, 3) = 5
```

Every hash has a model. The hash above: `a`, has the default model called *Object*.

A model is kind of like an interface, or a class. It has some properties that
the hash must contain, and defines some methods that are available to the hash.

You can define a model like this:

```go
vector := model (x, y);
```

The `vector` model is a model with two properties: *x* and *y*.

You can define methods on it with the following syntax:

```go
vector.print_something = fn (text) {
  print(this.x, text, this.y);
};
```

And then, you can instantiate the model:

```go
a := vector(3, 2);
b := vector(5, 1);

a.x = b.y;

a.print_something("and");
```
```shell
$ ./main

1 and 2
```

Of course, as the hash is still just a hash, you can add other properties to it.

You can also extend models. Here's a fairly basic example:

```go
animal := model (name, species);

animal.speak = fn () {
  print("Speak function not defined for", name);
};

dog := model (name) : animal (name, "dog");
fish := model (name) : animal (name, "fish");

dog.speak = fn () {
  print(name, "says 'Woof!'");
};
```

You can then instantiate animal, dog, or fish as a normal model. I think the syntax
for extending a model is really cool: `model (...args) : parent (...args)`. You
can probably figure out what it means. In the example above, the name of the dog
is set to the name argument given on instantiation, but the species is set automatically
to `"dog"`.

On a model, you can also define *special* methods, such as an initialization
method, or operator overloading. To demonstrate this, I'll go back to the vector
example:

```go
vector := model (x, y);

vector._new = fn () {
  print("new vector. x =", this.x, "y =", this.y);
  return this;
};

vector._plus = fn (other) {
  if (type(other) != vector) {
    err("expected another vector. got", type(other));
  };
  
  return vector(this.x + other.x, this.y + other.y);
};

vector._in = fn (other) {
  return other == this.x || other == this.y;
};

a := vector(2, 3);
print("a =", a);

b := vector(3, 2);
print("b =", b);

print("3 in b =", 3 in b);
print("a + b =", a + b);
```
```shell
$ ./main

new vector. x = 2 y = 3 
a = {x: 2, y: 3} 
new vector. x = 3 y = 2 
b = {x: 3, y: 2} 
3 in b = true 
new vector. x = 5 y = 5 
a + b = {x: 5, y: 5} 
```
