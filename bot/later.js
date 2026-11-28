// The minutes a reminder waits, rounded and held from 1 to 1440, and the
// message it sends then: the service, place, chat, user, name, email and
// answers of the message it was set in, its text the note and `later` the
// minutes.
const minutes = Math.min(Math.max(Math.round(Number(input.minutes)) || 1, 1), 1440)
const message = { text: String(input.note), later: minutes }
for (const k of ["service", "place", "chat", "user", "name", "email", "answers"]) {
  if (input.message[k] != null) message[k] = input.message[k]
}
return { minutes: minutes, message: message }
