import {Checkbox} from "primereact/checkbox";
import * as React from "react";
import {useEffect, useState} from "react";
import {IBoardSettings} from "../../components/Boards/IBoardSettings.ts";
import {getBoardByBoardSlug, getBoardSettings} from "../../components/Boards/BoardsService.ts";
import {Params, useParams} from "react-router";
import {InputTextarea} from "primereact/inputtextarea";
import {IBoardDetails} from "../../components/Boards/IBoardDetails.ts";
import {Button} from "primereact/button";
import {Card} from "primereact/card";

export function BoardSettingsPage() {
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string = params?.boardSlug || '';
    const [boardSettings, setBoardSettings] = useState<IBoardSettings>();
    const [boardDetails, setBoardDetails] = useState<IBoardDetails>();

    useEffect(() => {
        getBoardSettings(boardSlug).subscribe({
            next: (settings: IBoardSettings) => setBoardSettings(settings),
            error: (err) => console.error(err)
        });
        getBoardByBoardSlug(boardSlug).subscribe({
            next: (boardDetails: IBoardDetails) => setBoardDetails(boardDetails),
            error: (err) => console.error(err),
        })
    }, []);

    function updateSettings(settingName: string, value: string | boolean) {
        if (boardSettings) {
            setBoardSettings({...boardSettings, ...{[settingName]: value}});
        }
    }

    function updateBoardDetails(settingName: string, value: string | boolean) {
        if (boardDetails) {
            setBoardDetails({...boardDetails.board, ...{[settingName]: value}});
        }
    }

    return (
        <>
            <div className="flex justify-between mb-4">
                <h1 className="text-3xl font-bold">Board Settings</h1>
                <Button label="Save Settings" icon="pi pi-check"/>
            </div>

            <section className="flex justify-between gap-4">
                <div className="w-1/2">
                    <Card>
                        <div className="mb-4 pt-4 flex gap-4">
                            <Checkbox inputId="isOfficialCheckbox"
                                      onChange={e => updateSettings('isOfficial', !!e.checked)}
                                      checked={!!boardSettings?.isOfficial}></Checkbox>
                            <div>
                                <label
                                    className="block mb-2 cursor-pointer"
                                    htmlFor="isOfficialCheckbox"><strong>Official Board</strong>
                                    <p className="text-gray-400">
                                        Marks this board as official, adding an icon to the board name. This board also
                                        appears
                                        in
                                        the
                                        unfiltered post list.
                                    </p>
                                </label>
                            </div>
                        </div>

                        <div className="mb-4 flex gap-4">
                            <Checkbox inputId="isPostApprovalRequiredCheckbox"
                                      onChange={e => updateSettings('isPostApprovalRequired', !!e.checked)}
                                      checked={!!boardSettings?.isPostApprovalRequired}></Checkbox>
                            <div>
                                <label
                                    className="block mb-2 cursor-pointer"
                                    htmlFor="isPostApprovalRequiredCheckbox"><strong>Post Approval Required</strong>
                                    <p className="text-gray-400">
                                        Requires all posts to be approved before they are public
                                    </p>
                                </label>
                            </div>
                        </div>

                        <div className="mb-4 flex gap-4">
                            <Checkbox inputId="isVisibleCheckbox"
                                      onChange={e => updateBoardDetails('visible', !!e.checked)}
                                      checked={!!boardDetails?.board.visible}></Checkbox>
                            <div>
                                <label
                                    className="block mb-2 cursor-pointer"
                                    htmlFor="isVisibleCheckbox"><strong>Visible</strong>
                                    <p className="text-gray-400">
                                        Controls whether the board is publicly visible. Moderators and admins can still
                                        see
                                        it.
                                    </p>
                                </label>
                            </div>
                        </div>
                    </Card>
                </div>

                <div className="w-1/2">
                    <Card>
                        <div className="pt-4">
                            <label
                                className="block mb-2 cursor-pointer"
                                htmlFor="boardDescriptionTextbox"><strong>Board Description</strong>
                                <p className="text-gray-400">
                                    Appears in board details
                                </p>
                            </label>
                            <InputTextarea value={boardDetails?.board.description}
                                           onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setValue(e.target.value)}
                                           rows={5} cols={30}/>
                        </div>
                    </Card>
                </div>
            </section>
        </>
    )
}