# Image Diffing tool

This is a dumb and simple tool for generating and/or applying differences between
a set of images. It is intended for correcting 1-pixel errors in GoRender output,
or adding detail that is difficult to represent in the voxel realm.

## Usage 

To generate diffs from an input and edited directory:

```
imagediff.exe -input example/input -output example/output -compare example/edited -generate
```

To apply generated diffs:

```
imagediff.exe -input example/input -output example/output -compare example/diff -apply
```

Images are matched based on filename (i.e file names must be the same between directory).
If an image does not have an edited version or diff it is skipped.

Files with different sizes/bounds will not be processed. Files with differing colour
depths will produce strange results - likely either all or none of the compared file
will be considered to be different.