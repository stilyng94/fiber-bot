import './App.css'
import Cart from '@/components/custom/Cart'
import { ItemCard } from '@/components/custom/ItemCard'
import foodImage from "@/assets/react.svg"
import { useEffect, useState } from 'react'

const foods = [
  { id: 1, title: "Pizza", price: 10.00, image: foodImage },
  { id: 2, title: "Burger", price: 10.00, image: foodImage },
  { id: 3, title: "Chips", price: 10.00, image: foodImage },
  { id: 4, title: "Sausage", price: 10.00, image: foodImage }
] satisfies { id: number, title: string, price: number, image: string }[]

const telegram = window.Telegram.WebApp;

function App() {
  const [cartItems, setCartItems] = useState<{ id: number, quantity: number, price: number }[]>([])

  useEffect(() => {
    telegram.ready();

    return () => {
      telegram.close();
    }
  }, [])


  const onAddToCart = (food: {
    id: number;
    title: string;
    price: number;
    image: string;
  }) => {
    const exist = cartItems.findIndex((f) => f.id === food.id)
    if (exist >= 0) {
      setCartItems((prev) => {
        const data = [...prev]
        data[exist].quantity += 1
        return data
      })
    } else {
      setCartItems((prev) => [...prev, { id: food.id, quantity: 1, price: food.price }])
    }
  }

  const onRemoveFromCart = (food: {
    id: number;
    title: string;
    price: number;
    image: string;
  }) => {
    const exist = cartItems.findIndex((f) => f.id === food.id)
    if (exist >= 0 && cartItems[exist].quantity === 1) {
      setCartItems((prev) => {
        const data = prev.filter((f) => f.id !== food.id)
        return data
      })
    } else {
      setCartItems((prev) => {
        const data = [...prev]
        data[exist].quantity -= 1
        return data
      })
    }


  }

  const onCheckout = () => {
    telegram.MainButton.setText("Pay :)");
    telegram.MainButton.show();
  }


  return (
    <div className="container">
      <h1 className='text-center text-2xl font-bold'>Order Food</h1>
      <Cart cartItems={cartItems} onCheckout={onCheckout} />
      <div className="flex flex-wrap gap-5 mt-5 justify-center items-center w-full mx-auto px-4 sm:px-6 lg:px-8">
        {foods.map((food) => <ItemCard key={food.id} food={food} onAddToCart={onAddToCart} onRemoveFromCart={onRemoveFromCart} className='mb-5' />
        )}
      </div>

    </div>

  )
}

export default App
