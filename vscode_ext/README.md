# Sibylgo VSCode extension

## See VSCode extension errors

When packaging, you won't see compile errors, so simply run `npm run compile` to see them.

## Log something from VSCode extension

```javascript
    let orange = vscode.window.createOutputChannel("Orange");

    //Write to output.
    orange.appendLine("I am a banana.");
```

## Add new npm dependency to consume in extension

* npm install --save dep
* reference in extension
* exclude just the node_modules files you need in .vscodeignore, so it's actually added in package. E.g.: !node_modules/jquery/dist/jquery.min.js
