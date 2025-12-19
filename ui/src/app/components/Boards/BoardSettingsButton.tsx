import {useEffect, useState} from "react";
import {Button} from "primereact/button";

export function BoardSettingsButton() {
    const [isSettingsAreaAvailable, setIsSettingsAreaAvailable] = useState<boolean>(false);

    useEffect(() => {

    }, [])

    return (
        <>
            {isSettingsAreaAvailable && (
                <Button
                    label="Settings"
                    icon="pi pi-cog"
                    className="p-button-rounded"
                    onClick={() => goToSettingsArea(false)}
                />
            )}
        </>
    )
}