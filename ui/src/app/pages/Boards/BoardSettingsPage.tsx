import {Checkbox} from "primereact/checkbox";
import * as React from "react";
import {RefObject, useEffect, useRef, useState} from "react";
import {getBoardByBoardSlug, saveBoardDetails} from "../../components/Boards/BoardsService.ts";
import {NavigateFunction, Params, useNavigate, useParams} from "react-router";
import {InputTextarea} from "primereact/inputtextarea";
import {IBoardDetails} from "../../components/Boards/types/IBoardDetails.ts";
import {Button} from "primereact/button";
import {Card} from "primereact/card";
import {FileUpload} from "primereact/fileupload";
import {DumpsterFireError} from "../../components/Shared/DumpsterFireError.tsx";
import {Subscription} from "rxjs";
import {Toast} from "primereact/toast";
import {IBoardDetailsPayload} from "../../components/Boards/types/IBoardDetailsPayload.ts";

export function BoardSettingsPage() {
    const toast: RefObject<Toast | null> = useRef<Toast>(null);
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string = params?.boardSlug || '';
    const [boardDetails, setBoardDetails] = useState<IBoardDetails>();
    const [somethingWentWrong, setSomethingWentWrong] = useState<boolean>(false);
    const navigate: NavigateFunction = useNavigate();
    const saveSettings$ = React.useRef<Subscription>(null);

    useEffect(() => {
        getBoardByBoardSlug(boardSlug).subscribe({
            next: (boardDetails: IBoardDetails) => {
                setBoardDetails(boardDetails)
            },
            error: (err) => {
                console.error(err);
                setSomethingWentWrong(true);
            },
        })
    }, []);

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

    const onSaveButtonClicked = () => {
        if (boardDetails) {
            const payload: IBoardDetailsPayload = {
                isPrivate: boardDetails.board.isPrivate,
                isVisible: boardDetails.board.isVisible,
                isPostApprovalRequired: boardDetails.board.isPostApprovalRequired,
                isOfficial: boardDetails.board.isOfficial,
                description: boardDetails.board.description,
            }
            saveSettings$.current = saveBoardDetails(payload).subscribe({
                next: () => {
                    if (toast.current) {
                        toast.current.show({
                            severity: 'success',
                            summary: 'Success',
                            detail: 'Board settings saved',
                            life: 3000
                        });
                    }
                },
                error: (err) => {
                    console.error("Error saving board settings:", err);
                }
            })
        }
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
                            <Button label="Save Settings" icon="pi pi-check" onClick={onSaveButtonClicked}/>
                        </div>
                    </div>

                    <section className="flex justify-between gap-4">
                        <div className="w-1/2">
                            <Card>
                                <div className="mb-4 pt-4 flex gap-4">
                                    <div className="">
                                        <Checkbox inputId="isOfficialCheckbox"
                                                  onChange={e => updateBoardDetails('isOfficial', !!e.checked)}
                                                  checked={!!boardDetails?.board.isOfficial}></Checkbox>
                                    </div>
                                    <div>
                                        <label
                                            className="block mb-2 cursor-pointer"
                                            htmlFor="isOfficialCheckbox">
                                            <i className="pi pi-verified pr-1"/> <strong>Official Board</strong>
                                            <p className="text-gray-400">
                                                Marks this board as official, adding an icon to the board name. This
                                                board also appears in the unfiltered post list.
                                            </p>
                                        </label>
                                    </div>
                                </div>

                                <div className="mb-4 flex gap-4">
                                    <div className="">
                                        <Checkbox inputId="isPostApprovalRequiredCheckbox"
                                                  onChange={e => updateBoardDetails('isPostApprovalRequired', !!e.checked)}
                                                  checked={!!boardDetails?.board.isPostApprovalRequired}></Checkbox>
                                    </div>
                                    <div>
                                        <label
                                            className="block mb-2 cursor-pointer"
                                            htmlFor="isPostApprovalRequiredCheckbox">
                                            <i className="pi pi-thumbs-up-fill pr-1"/> <strong>Post Approval
                                            Required</strong>
                                            <p className="text-gray-400">
                                                Requires all posts from users to be approved before they are public.
                                                Moderators and admins can post without approval.
                                            </p>
                                        </label>
                                    </div>
                                </div>

                                <div className="mb-4 flex gap-4">
                                    <div className="">
                                        <Checkbox inputId="isVisibleCheckbox"
                                                  onChange={e => updateBoardDetails('isVisible', !!e.checked)}
                                                  checked={!!boardDetails?.board.isVisible}></Checkbox>
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
                                                  onChange={e => updateBoardDetails('isPrivate', !!e.checked)}
                                                  checked={!!boardDetails?.board.isPrivate}></Checkbox>
                                    </div>
                                    <div>
                                        <label
                                            className="block mb-2 cursor-pointer"
                                            htmlFor="isPrivateCheckbox">
                                            <i className="pi pi-users pr-1"/> <strong>Private</strong>
                                            <p className="text-gray-400">
                                                When the board is private, all members that join must be approved. Posts
                                                will only be visible to moderators and admins.
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
                                                           onChange={
                                                               (e: React.ChangeEvent<HTMLTextAreaElement>) => {
                                                                   updateBoardDetails('description', e.target.value)
                                                               }
                                                           }
                                                           rows={5} cols={30}/>
                                        </div>
                                    </div>
                                </section>
                            </Card>
                        </div>
                    </section>
                    <Toast ref={toast}/>
                </>
            )}
        </>
    )
}