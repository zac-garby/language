animal := model (name, type);
animal._new = fn () { print("Creating a new", this.type + "!"); };

sheep := model (name) : animal (name, "sheep");
sheep._new = fn () { print("New sheep :)"); };

pig := model (name) : animal (name, "pig");

s := sheep("x");
p := pig("y");