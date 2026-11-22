const assert = require("node:assert/strict");
const test = require("node:test");
const { parseCommandLine } = require("./args");

test("parses unquoted Windows paths without consuming backslashes", () => {
  assert.deepEqual(parseCommandLine("forge .\\Bitfile"), ["forge", ".\\Bitfile"]);
});

test("preserves spaces in quoted arguments", () => {
  assert.deepEqual(parseCommandLine('pack "C:\\Projects\\my app\\Bitfile"'), [
    "pack",
    "C:\\Projects\\my app\\Bitfile",
  ]);
});

test("rejects unmatched quotes", () => {
  assert.throws(() => parseCommandLine('forge "my Bitfile'), /Unclosed quote/);
});
