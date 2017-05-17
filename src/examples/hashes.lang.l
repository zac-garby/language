hash := {
  x: 5,
  y: 3,
  add: \(x, y) = x + y
};

hash.x; # 5
hash.y; # 3
hash.add(hash.x, hash.y); # 8