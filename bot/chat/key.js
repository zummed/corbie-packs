// The question's key, 12 random letters and digits, and the full location of
// its ballot, from which it takes close, or null with no ballot.
const digits = "0123456789abcdefghijklmnopqrstuvwxyz"
const key = Array.from({ length: 12 }, () => digits[Math.floor(Math.random() * digits.length)]).join("")
return { key: key, ballot: input.ballot ? input.tree + "/" + input.ballot.at : null }
