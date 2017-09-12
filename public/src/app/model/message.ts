export interface Contact {
  type: number,
  ips: string[],
  callsign: string
}

export interface Message {
  sequence: number,
  contact: Contact,
  message: string,
  ack: boolean
}