function parseCommandLine(input) {
  const args = [];
  let value = "";
  let quote = "";
  let started = false;

  for (let index = 0; index < input.length; index += 1) {
    const character = input[index];
    if (
      character === "\\" &&
      index + 1 < input.length &&
      (input[index + 1] === "\\" || input[index + 1] === quote || (!quote && (input[index + 1] === '"' || input[index + 1] === "'")))
    ) {
      value += input[index + 1];
      index += 1;
      started = true;
      continue;
    }
    if (quote) {
      if (character === quote) {
        quote = "";
      } else {
        value += character;
      }
      started = true;
      continue;
    }
    if (character === '"' || character === "'") {
      quote = character;
      started = true;
    } else if (/\s/.test(character)) {
      if (started) {
        args.push(value);
        value = "";
        started = false;
      }
    } else {
      value += character;
      started = true;
    }
  }

  if (quote) {
    throw new Error("Unclosed quote in Bitpack command.");
  }
  if (started) {
    args.push(value);
  }
  return args;
}

module.exports = { parseCommandLine };
