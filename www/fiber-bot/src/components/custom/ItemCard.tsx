import { Plus, Minus } from "lucide-react"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { useState } from "react"




type CardProps = React.ComponentProps<typeof Card> & {
  food: {
    id: number;
    title: string, price: number, image: string
  }, onAddToCart: (food: {
    id: number;
    title: string;
    price: number;
    image: string;
  }) => void,
  onRemoveFromCart: (food: {
    id: number;
    title: string;
    price: number;
    image: string;
  }) => void
}


export function ItemCard({ className, food, onAddToCart, onRemoveFromCart, ...props }: CardProps) {
  const [count, setCount] = useState(0)

  return (
    <Card className={cn("w-80 relative", className)} {...props}>
      {count > 0 ? <span className="absolute top-2 right-2 text-xs text-white bg-red-500 rounded-[50%] w-8 h-8 flex items-center justify-center transition-all scale-105 ease-in">{count}</span> : <></>}
      <CardContent className="grid gap-4">
        <div className="aspect-video pt-4">
          <img src={food.image} alt={food.title} className="object-contain w-full h-full" loading="lazy" decoding="async" />
        </div>
        <CardHeader className="flex-row items-baseline justify-center gap-2">
          <CardTitle className="font-normal">
            {food.title}
          </CardTitle>
          <CardTitle className="font-bold">${food.price}</CardTitle>
        </CardHeader>
      </CardContent>
      <CardFooter className="gap-2">
        <Button className="w-full" onClick={() => {
          setCount((prev) => prev += 1)
          onAddToCart(food)
        }}>
          <Plus className="mr-2 h-4 w-4" />
        </Button>
        {count > 0 ? <Button className="w-full" variant={"destructive"} onClick={() => {
          setCount((prev) => prev -= 1)
          onRemoveFromCart(food)
        }}>
          <Minus className="mr-2 h-4 w-4" />
        </Button> : <></>}
      </CardFooter>
    </Card>
  )
}
