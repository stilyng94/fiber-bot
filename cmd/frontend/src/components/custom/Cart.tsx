

type CartProps = {
  cartItems: {
    id: number, quantity: number, price: number
  }[],
}

function Cart({ cartItems }: CartProps) {
  const totalPrice = cartItems.reduce((prev, curr) => prev + (curr.quantity * curr.price), 0)

  return (
    <div className='mt-5 flex items-center justify-center flex-col'>
      <p>{cartItems.length === 0 ? "No items in cart" : ""}</p>
      <br />
      <div className="">
        <span className="font-bold mr-2">Total Price: {totalPrice.toFixed(2)}</span>
      </div>
    </div>
  )

}

export default Cart
