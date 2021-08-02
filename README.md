# ssn2svg

A command line tool for converting from SmallSketch Note file (.ssn) to SVG. 
About SmallSketch Note, please see this site: https://www.smallsketch.app/

![Coelacanth SVG](https://github.com/ssnote/ssn2svg/blob/main/examples/coelacanth.svg)


## Build

In order to create a ssn2svg command:

```
go build
```


## Usage

This coelacanth.ssn file is just an example. You can use your own ssn file instead of it.

```
cat examples/coelacanth.ssn | ssn2svg > coelacanth.svg
```

You can also apply style with style.json when converting.

```
cat examples/coelacanth.ssn | ssn2svg examples/style.json > coelacanth.svg
```

