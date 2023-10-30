/* eslint-disable @typescript-eslint/no-explicit-any */
import Cart from '@/components/custom/Cart'
import { ItemCard } from '@/components/custom/ItemCard'
import { useEffect, useState } from 'react'
import { addToCart, foods, validateAndUpsertUser } from '@/lib/api';
import { useLoaderData, useNavigate } from "react-router-dom";


const telegram = window.Telegram.WebApp;

function Homepage() {
  const navigate = useNavigate()
  telegram.BackButton.hide()
  useLoaderData()
  const [cartItems, setCartItems] = useState<{ id: number, quantity: number, price: number }[]>([]);

  useEffect(() => {
    const onCheckout = async () => {
      await addToCart()
      navigate({ pathname: "/checkout" })
    }
    telegram.onEvent("mainButtonClicked", onCheckout)
    return () => {
      telegram.offEvent("mainButtonClicked", onCheckout)
    }
  }, [navigate])

  useEffect(() => {
    const toggleTelegramMainButton = () => {
      if (cartItems.length === 0) {
        telegram.MainButton.hide();
      } else {
        telegram.MainButton.setText("To order(s)");
        telegram.MainButton.show();
      }
    }
    toggleTelegramMainButton();
  }, [cartItems])

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

  return (
    <div>
      <h1 className='text-center text-2xl font-bold'>Order Food</h1>
      <Cart cartItems={cartItems} />
      <div className="flex flex-wrap gap-5 mt-5 justify-center items-center w-full mx-auto px-4 sm:px-6 lg:px-8">
        {foods.map((food) => <ItemCard key={food.id} food={food} onAddToCart={onAddToCart} onRemoveFromCart={onRemoveFromCart} className='mb-5' />
        )}
      </div>
    </div>
  )
}

export default Homepage

export function loader() {
  return validateAndUpsertUser(telegram.initData);
}
