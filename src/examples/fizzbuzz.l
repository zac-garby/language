# fizz buzz program
# prints fizz if n is divisible by 3, buzz if it is divisible by 5, and fizzbuzz if both

for (i | 1..100) {
  n := i + 1;
  
  if (n % 3 == 0 && n % 5 == 0) {
    print("fizzbuzz");
  } elif (n % 3 == 0) {
    print("fizz");
  } elif (n % 5 == 0) {
    print("buzz");
  } else {
    print(n);
  };
};