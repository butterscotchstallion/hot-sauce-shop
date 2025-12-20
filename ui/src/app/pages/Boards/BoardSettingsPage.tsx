import {Checkbox} from "primereact/checkbox";
import * as React from "react";
import {useEffect, useState} from "react";
import {IBoardSettings} from "../../components/Boards/IBoardSettings.ts";
import {getBoardByBoardSlug, getBoardSettings} from "../../components/Boards/BoardsService.ts";
import {Params, useNavigate, useParams} from "react-router";
import {InputTextarea} from "primereact/inputtextarea";
import {IBoardDetails} from "../../components/Boards/IBoardDetails.ts";
import {Button} from "primereact/button";
import {Card} from "primereact/card";
import {FileUpload} from "primereact/fileupload";
import {DumpsterFireError} from "../../components/Shared/DumpsterFireError.tsx";

export function BoardSettingsPage() {
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string = params?.boardSlug || '';
    const [boardSettings, setBoardSettings] = useState<IBoardSettings>();
    const [boardDetails, setBoardDetails] = useState<IBoardDetails>();
    const [somethingWentWrong, setSomethingWentWrong] = useState<boolean>(false);
    const navigate = useNavigate();

    useEffect(() => {
        getBoardSettings(boardSlug).subscribe({
            next: (settings: IBoardSettings) => setBoardSettings(settings),
            error: (err) => {
                console.error(err);
                setSomethingWentWrong(true);
            }
        });
        getBoardByBoardSlug(boardSlug).subscribe({
            next: (boardDetails: IBoardDetails) => setBoardDetails(boardDetails),
            error: (err) => {
                console.error(err);
                setSomethingWentWrong(true);
            },
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

    function onUpload() {

    }

    const goToBoardPage = () => {
        navigate(`/boards/${boardSlug}`);
    }

    return (
        <>
            {somethingWentWrong ? (
                <DumpsterFireError/>
            ) : (
                <>
                    <div className="flex justify-between mb-4 gap-4">
                        <h1 className="text-3xl font-bold">Board Settings</h1>
                        <div className="w-3/4 gap-4 flex justify-end">
                            <Button label="View Board" icon="pi pi-eye" onClick={() => goToBoardPage()}/>
                            <Button label="Save Settings" icon="pi pi-check"/>
                        </div>
                    </div>

                    <section className="flex justify-between gap-4">
                        <div className="w-1/2">
                            <Card>
                                <div className="mb-4 pt-4 flex gap-4">
                                    <div className="">
                                        <Checkbox inputId="isOfficialCheckbox"
                                                  onChange={e => updateSettings('isOfficial', !!e.checked)}
                                                  checked={!!boardSettings?.isOfficial}></Checkbox>
                                    </div>
                                    <div>
                                        <label
                                            className="block mb-2 cursor-pointer"
                                            htmlFor="isOfficialCheckbox">
                                            <i className="pi pi-verified pr-1"/> <strong>Official Board</strong>
                                            <p className="text-gray-400">
                                                Marks this board as official, adding an icon to the board name. This
                                                board also
                                                appears
                                                in
                                                the
                                                unfiltered post list.
                                            </p>
                                        </label>
                                    </div>
                                </div>

                                <div className="mb-4 flex gap-4">
                                    <div className="">
                                        <Checkbox inputId="isPostApprovalRequiredCheckbox"
                                                  onChange={e => updateSettings('isPostApprovalRequired', !!e.checked)}
                                                  checked={!!boardSettings?.isPostApprovalRequired}></Checkbox>
                                    </div>
                                    <div>
                                        <label
                                            className="block mb-2 cursor-pointer"
                                            htmlFor="isPostApprovalRequiredCheckbox">
                                            <i className="pi pi-thumbs-up-fill pr-1"/> <strong>Post Approval
                                            Required</strong>
                                            <p className="text-gray-400">
                                                Requires all posts from users to be approved before they are public.
                                                Moderators
                                                and admins can post without approval.
                                            </p>
                                        </label>
                                    </div>
                                </div>

                                <div className="mb-4 flex gap-4">
                                    <div className="">
                                        <Checkbox inputId="isVisibleCheckbox"
                                                  onChange={e => updateBoardDetails('visible', !!e.checked)}
                                                  checked={!!boardDetails?.board.visible}></Checkbox>
                                    </div>
                                    <div>
                                        <label
                                            className="block mb-2 cursor-pointer"
                                            htmlFor="isVisibleCheckbox">
                                            <i className="pi pi-eye pr-1"/> <strong>Visible</strong>
                                            <p className="text-gray-400">
                                                Controls whether the board is publicly visible. Moderators and admins
                                                can still
                                                see it.
                                            </p>
                                        </label>
                                    </div>
                                </div>

                                <div className="mb-4 flex gap-4">
                                    <div className="">
                                        <Checkbox inputId="isPrivateCheckbox"
                                                  onChange={e => updateBoardDetails('private', !!e.checked)}
                                                  checked={!!boardDetails?.board.private}></Checkbox>
                                    </div>
                                    <div>
                                        <label
                                            className="block mb-2 cursor-pointer"
                                            htmlFor="isPrivateCheckbox">
                                            <i className="pi pi-users pr-1"/> <strong>Private</strong>
                                            <p className="text-gray-400">
                                                When the board is private, all members that join must be approved. Posts
                                                will
                                                only
                                                be visible to moderators and admins.
                                            </p>
                                        </label>
                                    </div>
                                </div>
                            </Card>
                        </div>

                        {/* Board details */}
                        <div className="w-1/2">
                            <Card>
                                <section className="flex justify-between">
                                    <div className="w-[128px] h-[128px] mt-4">
                                        <img src="/images/hot-pepper.png" width={128} height={128}
                                             alt="Board Thumbnail"/>
                                    </div>
                                    <div className="w-3/4">
                                        <div className="pt-4 mb-4">
                                            <label
                                                className="block mb-2 cursor-pointer"
                                                htmlFor="boardThumbnailFilename">
                                                <i className="pi pi-image pr-1"/> <strong>Board Thumbnail</strong>
                                            </label>
                                            <FileUpload
                                                mode="basic"
                                                name="boardThumbnailFilename"
                                                url={``}
                                                accept="image/*"
                                                maxFileSize={1000000}
                                                onUpload={onUpload}
                                            />
                                        </div>
                                        <div className="">
                                            <label
                                                className="block mb-2 cursor-pointer"
                                                htmlFor="boardDescriptionTextbox">
                                                <i className="pi pi-file pr-1"/> <strong>Board Description</strong>
                                            </label>
                                            <InputTextarea value={boardDetails?.board.description}
                                                           onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setValue(e.target.value)}
                                                           rows={5} cols={30}/>
                                        </div>
                                    </div>
                                </section>
                            </Card>
                        </div>
                    </section>
                </>
            )}
        </>
    )
}