import './App.scss'
import {PrimeReactProvider} from "primereact/api";

function App() {
    return (
        <PrimeReactProvider>
            <h1 className="text-3xl font-bold underline">Hot Sauce Shop</h1>
        </PrimeReactProvider>
    )
}

export default App
