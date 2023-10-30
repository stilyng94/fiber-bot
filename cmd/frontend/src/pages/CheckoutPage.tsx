/* eslint-disable @typescript-eslint/no-explicit-any */
import { Button } from "@/components/ui/button"
import { useEffect } from "react";
import { useLoaderData, useNavigate } from "react-router-dom";
import foodImage from "@/assets/react.svg";
import { checkout } from "@/lib/api";


export const cartItems = [
  { id: 1, image: foodImage, price: 10, title: "burger" },
] satisfies {
  id: number;
  title: string;
  price: number;
  image: string;
}[];



const telegram = window.Telegram.WebApp;

function CheckoutPage() {
  const navigate = useNavigate()
  const data = useLoaderData() as any


  useEffect(() => {
    const invoiceHandler = () => {
      telegram.openInvoice(data.invoice, async (status) => {
        if (status === 'paid') {
          return navigate("/");
        }
      });
    }
    const backHandler = () => navigate(-1);

    telegram.BackButton.show()
    telegram.MainButton.setText(`Pay $${data.totalCharge}`);
    telegram.MainButton.show();
    telegram.MainButton.onClick(invoiceHandler);
    telegram.BackButton.onClick(backHandler);
    return () => {
      telegram.BackButton.offClick(backHandler);
      telegram.MainButton.offClick(invoiceHandler);
    }
  }, [data, navigate])



  return (
    <div className="flex flex-col gap-4 mt-4" >
      <div className="flex flex-row justify-between items-center">
        <h1 className='text-center text-2xl font-bold'>Your Orders</h1>
        <Button variant={"outline"}>
          Edit
        </Button>
      </div>
      {cartItems.map((item) => (
        <div className="flex flex-row justify-between items-center">
          <div className="flex flex-row gap-2" key={item.id}>
            <img src={item.image} alt="image" className="object-contain w-10 h-10" loading="lazy" decoding="async" />
            <div className="flex flex-col">
              <h1 className='text-center font-bold' >{`${item.title} ${2}X`}</h1>
              <h1 className='text-center font-light'>Chips</h1>
            </div>
          </div>
          <h1 className='text-center text-2xl font-bold'>${item.price.toFixed(2)}</h1>
        </div>
      ))}
    </div>
  )
}

export default CheckoutPage


export function loader() {
  telegram.BackButton.hide()
  return checkout();
}
