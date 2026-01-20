import {NavigateFunction, useNavigate} from "react-router";

interface IBoardSettingsButtonProps {
    settingsAreaAvailable: boolean;
    boardSlug: string;
}

export function BoardSettingsButton({settingsAreaAvailable, boardSlug}: IBoardSettingsButtonProps) {
    const navigate: NavigateFunction = useNavigate();

    const goToSettingsArea = () => {
        if (settingsAreaAvailable) {
            navigate(`/boards/${boardSlug}/settings`);
        }
    }

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