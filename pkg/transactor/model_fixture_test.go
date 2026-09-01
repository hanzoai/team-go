package transactor

// testModel is a small platform model: the classes these tests actually exercise.
//
// The real artifact is not in this repository and is not compiled in (model.go),
// so the tests carry their own. Naming the protocol's identifiers is what talking
// to it requires — the same reason the Go in this package names them — and a
// hand-written handful is not a copy of anybody's bundle.
//
// ClassifierKind: CLASS=0, INTERFACE=1, MIXIN=2 (hierarchy.go).
var testModel = []byte(`[
  {"_class":"core:class:TxCreateDoc","objectClass":"core:class:Class",
   "objectId":"core:class:Doc","attributes":{"kind":0}},
  {"_class":"core:class:TxCreateDoc","objectClass":"core:class:Class",
   "objectId":"core:class:Space","attributes":{"kind":0,"extends":"core:class:Doc"}},
  {"_class":"core:class:TxCreateDoc","objectClass":"core:class:Class",
   "objectId":"core:class:SystemSpace","attributes":{"kind":0,"extends":"core:class:Space"}},
  {"_class":"core:class:TxCreateDoc","objectClass":"core:class:Class",
   "objectId":"contact:class:Person","attributes":{"kind":0,"extends":"core:class:Doc"}},
  {"_class":"core:class:TxCreateDoc","objectClass":"core:class:Class",
   "objectId":"contact:class:SocialIdentity","attributes":{"kind":0,"extends":"core:class:Doc"}},
  {"_class":"core:class:TxCreateDoc","objectClass":"core:class:Class",
   "objectId":"contact:class:PersonSpace","attributes":{"kind":0,"extends":"core:class:Space"}},
  {"_class":"core:class:TxCreateDoc","objectClass":"core:class:Class",
   "objectId":"contact:mixin:Employee","attributes":{"kind":2,"extends":"contact:class:Person"}},
  {"_class":"core:class:TxCreateDoc","objectClass":"core:class:Class",
   "objectId":"chunter:class:Channel","attributes":{"kind":0,"extends":"core:class:Space"}},
  {"_class":"core:class:TxCreateDoc","objectClass":"core:class:Class",
   "objectId":"chunter:class:ChatMessage","attributes":{"kind":0,"extends":"core:class:Doc"}},
  {"_class":"core:class:TxCreateDoc","objectClass":"core:class:Class",
   "objectId":"tracker:class:Project","attributes":{"kind":0,"extends":"core:class:Space"}}
]`)
