import { Button } from "@/components/ui/button";

type CartProps = {
  cartItems: {
    id: number, quantity: number, price: number
  }[],
  onCheckout: () => void
}

function Cart({ cartItems, onCheckout }: CartProps) {
  const totalPrice = cartItems.reduce((prev, curr) => prev + (curr.quantity * curr.price), 0)

  return (
    <div className='mt-5'>
      <p>{cartItems.length === 0 ? "No items in cart" : ""}</p>
      <br />
      <span className="font-bold mr-2">Total Price: {totalPrice.toFixed(2)}</span>
      <Button disabled={cartItems.length === 0} aria-disabled={cartItems.length === 0} onClick={onCheckout}>{cartItems.length === 0 ? "Order Now" : "Checkout"}</Button>
    </div>
  )
}

export default Cart
