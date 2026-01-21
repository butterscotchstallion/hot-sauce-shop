import {Checkbox} from "primereact/checkbox";
import * as React from "react";
import {RefObject, useEffect, useMemo, useRef, useState} from "react";
import {
    deactivateBoard,
    getBoardByBoardSlug,
    isSettingsAreaAvailable,
    saveBoardDetails
} from "../../components/Boards/BoardsService.ts";
import {NavigateFunction, Params, useNavigate, useParams} from "react-router";
import {InputTextarea} from "primereact/inputtextarea";
import {IBoardDetails} from "../../components/Boards/types/IBoardDetails.ts";
import {Button} from "primereact/button";
import {Card} from "primereact/card";
import {FileUpload} from "primereact/fileupload";
import {DumpsterFireError} from "../../components/Shared/DumpsterFireError.tsx";
import {Subject, Subscription} from "rxjs";
import {Toast} from "primereact/toast";
import {IBoardDetailsPayload} from "../../components/Boards/types/IBoardDetailsPayload.ts";
import {InputNumber, InputNumberChangeEvent} from "primereact/inputnumber";
import {confirmPopup, ConfirmPopup} from "primereact/confirmpopup";
import {IUser} from "../../components/User/types/IUser.ts";
import {useSelector} from "react-redux";
import {RootState} from "../../store.ts";
import {isSuperMessageBoardAdmin} from "../../components/User/UserService.ts";
import {IUserRole} from "../../components/User/types/IUserRole.ts";
import {IBoard} from "../../components/Boards/types/IBoard.ts";
import {Message} from "primereact/message";
import {BoardNameLink} from "../../components/Boards/BoardNameLink.tsx";

export function BoardSettingsPage() {
    const user: IUser | null = useSelector((state: RootState) => state.user.user);
    const roles: IUserRole[] = useSelector((state: RootState) => state.user.roles);
    const toast: RefObject<Toast | null> = useRef<Toast>(null);
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string = params?.boardSlug || '';
    const [somethingWentWrong, setSomethingWentWrong] = useState<boolean>(false);
    const navigate: NavigateFunction = useNavigate();
    const saveSettings$ = React.useRef<Subscription>(null);
    const [boardUpdatePayload, setBoardUpdatePayload] = useState<IBoardDetailsPayload>();
    const [isBoardOwner, setIsBoardOwner] = useState<boolean>(false);
    const canDeactivateBoard = useMemo(() => {
        return isBoardOwner || isSuperMessageBoardAdmin(roles);
    }, [isBoardOwner, roles]);
    const [board, setBoard] = useState<IBoard>();
    const deactivateBoardSubjectRef = useRef<Subject<boolean | string>>(null);

    useEffect(() => {
        isSettingsAreaAvailable(boardSlug, roles).subscribe({
            next: (available: boolean) => {
                if (!available) {
                    setSomethingWentWrong(true);
                }
            },
            error: (err: string) => {
                console.error(err);
                setSomethingWentWrong(true);
            }
        });
        getBoardByBoardSlug(boardSlug).subscribe({
            next: (boardDetails: IBoardDetails) => {
                setBoard(boardDetails.board);
                setIsBoardOwner(user?.id === boardDetails.board.createdAtByUserId);
                setBoardUpdatePayload({
                    isPrivate: boardDetails.board.isPrivate,
                    isVisible: boardDetails.board.isVisible,
                    isPostApprovalRequired: boardDetails.board.isPostApprovalRequired,
                    isOfficial: boardDetails.board.isOfficial,
                    description: boardDetails.board.description,
                    thumbnailFilename: boardDetails.board.thumbnailFilename,
                    minKarmaRequiredToPost: boardDetails.board.minKarmaRequiredToPost,
                });
            },
            error: (err) => {
                console.error(err);
                setSomethingWentWrong(true);
            },
        });
        return () => {
            saveSettings$.current?.unsubscribe();
            deactivateBoardSubjectRef?.current?.unsubscribe();
        }
    }, []);

    function updateBoardDetails(settingName: string, value: string | boolean | number) {
        if (boardUpdatePayload) {
            setBoardUpdatePayload({...boardUpdatePayload, ...{[settingName]: value}});
        }
    }

    // TODO: set thumbnailFilename here, refactor endpoint to use a form
    // similar to how new post form works
    function onUpload() {

    }

    const goToBoardPage = () => {
        navigate(`/boards/${boardSlug}`);
    }

    const onSaveButtonClicked = () => {
        if (boardUpdatePayload) {
            saveSettings$.current = saveBoardDetails(boardSlug, boardUpdatePayload).subscribe({
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

    const confirmDeactivateBoard = (event) => {
        const accept = () => {
            deactivateBoardSubjectRef.current = deactivateBoard(boardSlug);
            deactivateBoardSubjectRef.current.subscribe({
                next: (deactivated: boolean) => {
                    if (toast.current) {
                        toast.current.show({
                            severity: 'success',
                            summary: 'Board deactivated',
                            detail: board?.displayName + ' has been deactivated successfully.',
                            life: 3000
                        });
                    }
                },
                error: (err: string) => {
                    console.error("Error deactivating board:", err);
                    if (toast.current) {
                        toast.current.show({
                            severity: 'error',
                            summary: 'Error',
                            detail: 'Error deactivating board: ' + err,
                            life: 3000,
                        })
                    }
                }
            });
        }
        const reject = () => {
        }
        confirmPopup({
            target: event.currentTarget,
            message: 'Are you sure you want to deactivate this board?',
            icon: 'pi pi-info-circle',
            defaultFocus: 'reject',
            acceptClassName: 'p-button-danger',
            accept,
            reject
        });
    };

    return (
        <>
            {somethingWentWrong ? (
                <DumpsterFireError/>
            ) : (
                <>
                    <div className="flex justify-between mb-4 gap-4">
                        <div className="w-1/2">
                            <h1 className="text-3xl font-bold">Board Settings</h1>
                            {board && (
                                <small className="italics">Viewing settings for
                                    <span className="ml-2">
                                        <BoardNameLink
                                            isOfficial={board.isOfficial}
                                            displayName={board.displayName}
                                            slug={board.slug}/>
                                    </span>
                                </small>
                            )}
                        </div>
                        <div className="w-1/2 gap-4 flex justify-end">
                            <Button label="View Board" icon="pi pi-eye" onClick={() => goToBoardPage()}/>
                            <Button label="Save Settings" icon="pi pi-check" onClick={onSaveButtonClicked}/>
                        </div>
                    </div>

                    {board?.deactivatedByUserId && board?.deactivatedAt && (
                        <section className="mb-2">
                            <Message className="w-full"
                                     severity="warn"
                                     text="This board is deactivated. It is not be visible to users."/>
                        </section>
                    )}

                    <section className="flex justify-between gap-4">
                        <div className="w-1/2">
                            <Card>
                                <div className="mb-4 pt-4 flex gap-4">
                                    <div className="">
                                        <Checkbox inputId="isOfficialCheckbox"
                                                  onChange={e => updateBoardDetails('isOfficial', !!e.checked)}
                                                  checked={!!boardUpdatePayload?.isOfficial}></Checkbox>
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
                                                  checked={!!boardUpdatePayload?.isPostApprovalRequired}></Checkbox>
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
                                                  checked={!!boardUpdatePayload?.isVisible}></Checkbox>
                                    </div>
                                    <div>
                                        <label
                                            className="block mb-2 cursor-pointer"
                                            htmlFor="isVisibleCheckbox">
                                            <i className="pi pi-eye pr-1"/> <strong>Visible</strong>
                                            <p className="text-gray-400">
                                                Controls whether the board is publicly visible. Moderators and admins
                                                can still see it.
                                            </p>
                                        </label>
                                    </div>
                                </div>

                                <div className="mb-4 flex gap-4">
                                    <div className="">
                                        <Checkbox inputId="isPrivateCheckbox"
                                                  onChange={e => updateBoardDetails('isPrivate', !!e.checked)}
                                                  checked={!!boardUpdatePayload?.isPrivate}></Checkbox>
                                    </div>
                                    <div>
                                        <label
                                            className="block mb-2 cursor-pointer"
                                            htmlFor="isPrivateCheckbox">
                                            <i className="pi pi-users pr-1"/> <strong>Private</strong>
                                            <p className="text-gray-400">
                                                When the board is private, all members must be approved prior to
                                                joining. Posts will only be visible to members.
                                            </p>
                                        </label>
                                    </div>
                                </div>

                                <div className="mb-4 flex gap-4">
                                    <div>
                                        <InputNumber
                                            onChange={(e: InputNumberChangeEvent) => {
                                                updateBoardDetails(
                                                    'minKarmaRequiredToPost',
                                                    parseInt(String((e.value || 0)), 10)
                                                )
                                            }}
                                            inputId="minKarmaRequiredToPost"
                                            value={boardUpdatePayload?.minKarmaRequiredToPost}
                                            min={0}
                                            max={1000000}
                                            maxLength={7}
                                            showButtons
                                        />
                                    </div>
                                    <div className="w-2/3">
                                        <label
                                            className="block mb-2 cursor-pointer"
                                            htmlFor="minKarmaRequiredToPost">
                                            <i className="pi pi-users pr-1"/> <strong>Minimum Karma Required To
                                            Post</strong>
                                            <p className="text-gray-400">
                                                Users will not be able to post to this board if their karma is below the
                                                value specified here.
                                            </p>
                                        </label>
                                    </div>
                                </div>
                            </Card>
                        </div>

                        {/* Board details */}
                        <div className="w-1/2">
                            <section className="mb-4">
                                <Card>
                                    <section className="flex justify-between">
                                        <div className="w-[128px] h-[128px] mt-4">
                                            <img src="/images/hot-pepper.png"
                                                 width={128}
                                                 height={128}
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
                                                <InputTextarea value={boardUpdatePayload?.description}
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
                            </section>
                            <section>
                                <h1 className="text-2xl font-bold mb-4">Danger Zone</h1>
                                <Card className="border-solid border-red-500 border-1">
                                    <Button
                                        disabled={!canDeactivateBoard}
                                        onClick={confirmDeactivateBoard}
                                        severity="danger"
                                        label="Deactivate Board"
                                        icon="pi pi-trash"/>
                                    <ConfirmPopup/>
                                </Card>
                            </section>
                        </div>
                    </section>
                </>
            )}
            <Toast ref={toast}/>
        </>
    )
}