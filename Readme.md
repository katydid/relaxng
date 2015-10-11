#Converts RelaxNG to Katydid

Katydid still requires some work to fully support XML.
When this is done testing of the translations can start.

## Known Issues

There are quite a few known issues:
  - Only simplified grammars are supported.
  - namespaces are not supported.
  - datatypes: only string and token are currently supported.
  - datatypeLibraries are not supported.

I don't really intend to fix these, but you never know.

### Only handles simplified relaxng grammars.

http://www.kohsuke.org/relaxng/rng2srng/ seems to be quite effective at converting the full spectrum of what is possible within the relaxng grammar to the simplified grammar.

```
java -jar rng2srng.jar full.rng > simplified.rng
```

