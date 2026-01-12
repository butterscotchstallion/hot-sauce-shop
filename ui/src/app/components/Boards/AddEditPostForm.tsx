import {IBoardPost} from "./types/IBoardPost.ts";
import {InputText} from "primereact/inputtext";
import {InputTextarea} from "primereact/inputtextarea";
import {Button} from "primereact/button";
import * as z from "zod";
import {ZodIssue} from "zod";
import {IFormErrata} from "../Shared/IFormErrata.ts";
import * as React from "react";
import {ChangeEvent, RefObject, useEffect, useRef, useState} from "react";
import {PostSchema} from "./PostSchema.ts";
import {Subject} from "rxjs";
import {addPost} from "./BoardsService.ts";
import {Toast} from "primereact/toast";
import {INewBoardPost} from "./types/INewBoardPost.ts";
import {NavigateFunction, useNavigate} from "react-router";
import {FileUpload, FileUploadSelectEvent} from "primereact/fileupload";
import "../../pages/Boards/NewPostPage.css";
import {useDebounce} from "primereact/hooks";

interface AddEditPostFormProps {
    post?: IBoardPost;
    boardSlug: string;
    parentSlug?: string;
    addPostCallback?: () => void;
}

interface IAddEditPostFormState {
    title?: string;
    postText: string;
}

export default function AddEditPostForm({post, boardSlug, parentSlug, addPostCallback}: AddEditPostFormProps) {
    let addPost$: Subject<IBoardPost>;
    const boardSlugRef: RefObject<string> = useRef<string>(boardSlug);
    const [isValid, setIsValid] = useState<boolean>(false);
    const toast: RefObject<Toast | null> = useRef<Toast>(null);
    const [postTitle, setPostTitle] = useState<string>("");
    const [postText, setPostText] = useState<string>("");
    const [debouncedPostText] = useDebounce(postText, 1000);
    const [debouncedPostTitle] = useDebounce(postTitle, 1000);
    const [postImages, setPostImages] = useState<File[]>([]);
    const navigate: NavigateFunction = useNavigate();
    const navigateToNewPost = (newPost: IBoardPost) => {
        let url = `/boards/${boardSlug}/posts/`;
        if (parentSlug) {
            url += parentSlug;
        } else {
            url += newPost.slug;
        }
        navigate(url);
    };
    const uploadOptions = {icon: '', iconOnly: true, className: 'hidden-upload-button'};
    const defaultErrata: IFormErrata = {
        name: '',
        postText: '',
    };
    const [formErrata, setFormErrata] = useState<IFormErrata>(defaultErrata);
    const resetErrata = () => {
        setFormErrata(defaultErrata);
        setIsValid(true);
    }
    const resetForm = () => {
        setPostText("");
        setPostTitle("");
        setFormErrata(defaultErrata);
        setIsValid(false);
    }
    const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        const post: INewBoardPost = {
            title: postTitle,
            postText: postText
        }
        if (parentSlug) {
            post.parentSlug = parentSlug;
        }
        const valid: boolean = validate();
        if (valid) {
            addPost$ = addPost(post, boardSlugRef.current, postImages);
            addPost$.subscribe({
                next: (newPost: IBoardPost) => {
                    if (toast.current) {
                        toast?.current.show({
                            severity: 'success',
                            summary: 'Success',
                            detail: 'Post added successfully',
                            life: 3000,
                        })
                    }
                    resetForm();
                    if (addPostCallback) {
                        addPostCallback();
                    }
                    navigateToNewPost(newPost);
                },
                error: (err) => {
                    console.log(err);
                    if (toast.current) {
                        toast?.current.show({
                            severity: 'error',
                            summary: 'Error',
                            detail: 'Error adding post: ' + err + '.',
                            life: 3000,
                        })
                    }
                }
            });
        } else {
            if (toast.current) {
                toast?.current.show({
                    severity: 'error',
                    summary: 'Error',
                    detail: 'Error adding post: ' + formErrata.name + '.',
                    life: 3000,
                })
            }
        }
    }
    const validate = (): boolean => {
        // TODO: fix this
        return true;
        console.log('validating');
        try {
            const fieldsToParse: IAddEditPostFormState = {postText};
            if (!parentSlug) {
                fieldsToParse.title = postTitle;
            }
            PostSchema.parse(fieldsToParse);
            resetErrata();
            return true;
        } catch (err: z.ZodError | unknown) {
            if (err instanceof z.ZodError) {
                const newErrata: IFormErrata = {...formErrata};
                err.issues.forEach((issue: ZodIssue) => {
                    newErrata[issue.path[0]] = issue.message;
                });
                setFormErrata(newErrata);
                setIsValid(false);
            }
            return false;
        }
    }

    const onFilesSelected = (e: FileUploadSelectEvent) => {
        setPostImages(e.files);
    }

    useEffect(() => {
        if (post) {
            setPostTitle(post.title);
            setPostText(post.postText);
        } else {
            setPostTitle("");
            setPostText("");
            resetErrata();
        }
        return () => {
            addPost$?.unsubscribe();
        }
    }, [post]);

    useEffect(() => {
        setIsValid(validate());
    }, [debouncedPostText, debouncedPostTitle]);

    return (
        <>
            <form onSubmit={onSubmit} className="w-full m-0 p-0">
                {!parentSlug && (
                    <>
                        <div className="w-full mb-4">
                            <label className="mb-2 block" htmlFor="post-title">Title</label>
                            <InputText
                                className="w-full"
                                onChange={(e) => {
                                    setPostTitle(e.target.value);
                                }}
                                value={postTitle}
                                maxLength={150}
                                invalid={!!formErrata.postTitle}
                                id="post-title"/>
                            <p className="text-red-500 pt-2">{formErrata.postTitle}</p>
                        </div>

                        <div className="w-full mb-4">
                            <FileUpload name="postImages[]"
                                        url={''}
                                        multiple
                                        accept="image/*"
                                        maxFileSize={10000000}
                                        customUpload={true}
                                        uploadOptions={uploadOptions}
                                        uploadHandler={e => console.log(e)}
                                        onSelect={e => onFilesSelected(e)}
                                        emptyTemplate={<p className="m-0">Drag and drop files to here to upload.</p>}/>
                        </div>
                    </>
                )}

                <div className="w-full">
                    <label className="mb-2 block" htmlFor="post-text">Post text</label>
                    <InputTextarea
                        className="w-full"
                        onChange={(e: ChangeEvent<HTMLTextAreaElement>) => {
                            setPostText(e.target.value);
                        }}
                        value={postText}
                        invalid={!!formErrata.postText}
                        rows={5}
                        cols={30}
                        id="post-text"/>
                    <p className="text-red-500 pt-2">{formErrata.postText}</p>
                </div>

                <div className="w-full flex mt-4 justify-end">
                    <Button type="submit" disabled={!isValid}><i className="pi pi-plus mr-2"></i> Post</Button>
                </div>
            </form>
            <Toast ref={toast}/>
        </>
    )
}