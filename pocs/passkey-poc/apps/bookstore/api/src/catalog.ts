export type Book = { id: string; title: string; author: string; year: number; price: number; tag: string; color: string };

export const BOOKS: Book[] = [
  { id: "b1", title: "Applied Cryptography", author: "Bruce Schneier", year: 1996, price: 5499, tag: "crypto", color: "#f4c7a1" },
  { id: "b2", title: "Serious Cryptography", author: "Jean-Philippe Aumasson", year: 2024, price: 3999, tag: "crypto", color: "#b8d8c8" },
  { id: "b3", title: "Security Engineering", author: "Ross Anderson", year: 2020, price: 6299, tag: "engineering", color: "#c9c3e6" },
  { id: "b4", title: "The Web Application Hacker's Handbook", author: "Stuttard and Pinto", year: 2011, price: 4499, tag: "web", color: "#f2b5b5" },
  { id: "b5", title: "Real-World Cryptography", author: "David Wong", year: 2021, price: 4299, tag: "crypto", color: "#a9cfe8" },
  { id: "b6", title: "Threat Modeling", author: "Adam Shostack", year: 2014, price: 3799, tag: "engineering", color: "#f6e3a1" },
  { id: "b7", title: "Hacking: The Art of Exploitation", author: "Jon Erickson", year: 2008, price: 3499, tag: "offense", color: "#d7e8b0" },
  { id: "b8", title: "Designing Secure Software", author: "Loren Kohnfelder", year: 2021, price: 3299, tag: "engineering", color: "#e8c3dc" },
];

export function findBook(id: string): Book | undefined {
  return BOOKS.find((b) => b.id === id);
}
