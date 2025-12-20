import {useEffect, useState} from "react";
import {NavigateFunction, Params, useNavigate, useParams} from "react-router";
import {isSettingsAreaAvailable} from "./BoardsService.ts";

export function BoardSettingsButton() {
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string = params?.boardSlug || '';
    const navigate: NavigateFunction = useNavigate();
    const [settingsAreaAvailable, setSettingsAreaAvailable] = useState<boolean>(false);

    const goToSettingsArea = () => {
        if (settingsAreaAvailable) {
            navigate(`/boards/${boardSlug}/settings`);
        }
    }

    useEffect(() => {
        isSettingsAreaAvailable(boardSlug).subscribe({
            next: (available: boolean) => setSettingsAreaAvailable(available),
            error: () => setSettingsAreaAvailable(false)
        })
    }, [])

    return (
        <>
            {settingsAreaAvailable && (
                <i
                    title="Board Settings"
                    className="pi pi-cog cursor-pointer hover:text-yellow-200"
                    onClick={() => goToSettingsArea()}
                />
            )}
        </>
    )
}