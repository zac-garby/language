# TODO list

 - Extend the standard library
 - Possible change the string type to an array of characters (like Haskell)
  - This would involve making a new object type (char), e.g. 'x' or 'y'
  - It would mean that strings can be more easily worked with
 - Add more builtin models
 - Convert all basic types to builtin models
 - Allow `model (..) : model` (without any parent args.) It should just pass all
  model args to the parent