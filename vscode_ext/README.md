# Sibylgo VSCode extension

## See VSCode extension errors

When packaging, you won't see compile errors, so simply run `npm run compile` to see them.

## Log something from VSCode extension

```javascript
    let orange = vscode.window.createOutputChannel("Orange");

    //Write to output.
    orange.appendLine("I am a banana.");
```
