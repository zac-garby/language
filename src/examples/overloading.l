vec := model (x, y);

vec._new = fn () {
  print("new vector. x =", this.x, "y =", this.y);
};

vec._plus = fn (other) {
  if (type(other) != vec) {
    err("expected another vector. got", type(other));
  };
  
  return vec(this.x + other.x, this.y + other.y);
};

vec._in = fn (other) {
  return other == this.x || other == this.y;
};

a := vec(2, 3);
b := vec(3, 2);
c := a + b;

print(c); # {x: 5, y: 5}
print(3 in a); # true
