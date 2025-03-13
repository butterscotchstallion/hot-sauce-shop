import {Toast, ToastMessageOptions} from "primereact/toast";
import {Context, createContext, ReactElement, RefObject, useContext, useRef} from "react";

const ToastContext: Context<undefined> = createContext(undefined);

export const ToastContextProvider = ({children}): ReactElement => {
    const toastRef: RefObject<null> = useRef(null);

    const showToast = (options: ToastMessageOptions) => {
        if (!toastRef.current) return;
        toastRef.current.show(options);
    };

    return (
        <ToastContext.Provider value={{showToast}}>
            <Toast ref={toastRef}/>
            <div>{children}</div>
        </ToastContext.Provider>
    );
};

export const useToastContext = () => {
    const context = useContext(ToastContext);

    if (!context) {
        throw new Error(
            "useToastContext have to be used within ToastContextProvider"
        );
    }

    return context;
};
