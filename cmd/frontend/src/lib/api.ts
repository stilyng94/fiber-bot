/* eslint-disable @typescript-eslint/no-explicit-any */
import foodImage from "@/assets/react.svg";

export const items = [
  { productId: 1, title: "Pizza", quantity: 3 },
  { productId: 2, title: "Burger", quantity: 3 },
  { productId: 3, title: "Chips", quantity: 3 },
  { productId: 4, title: "Sausage", quantity: 3 },
];
export const foods = [
  { id: 1, title: "Pizza", price: 10.0, image: foodImage },
  { id: 2, title: "Burger", price: 10.0, image: foodImage },
  { id: 3, title: "Chips", price: 10.0, image: foodImage },
  { id: 4, title: "Sausage", price: 10.0, image: foodImage },
] satisfies { id: number; title: string; price: number; image: string }[];

export const validateAndUpsertUser = async (token: string) => {
  const response = await fetch("/users", {
    body: JSON.stringify({ token }),
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    mode: "cors",
  });
  if (!response.ok) {
    throw { message: "Request error" };
  }
  return response.json();
};

export const checkout = async (): Promise<any> => {
  const response = await fetch("/orders", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    mode: "cors",
  });
  if (!response.ok) {
    throw { message: "request error" };
  }
  return response.json();
};

export const addToCart = async (): Promise<any> => {
  const response = await fetch("/orders/cart", {
    body: JSON.stringify({ items }),
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    mode: "cors",
  });
  if (!response.ok) {
    throw { message: "request error" };
  }
  return response.json();
};
