import {IBoardPost} from "./IBoardPost.ts";
import {InputText} from "primereact/inputtext";
import {InputTextarea} from "primereact/inputtextarea";
import {Button} from "primereact/button";
import * as z from "zod";
import {ZodIssue} from "zod";
import {IFormErrata} from "../Shared/IFormErrata.ts";
import * as React from "react";
import {RefObject, useEffect, useRef, useState} from "react";
import {PostSchema} from "./PostSchema.ts";
import {Subject} from "rxjs";
import {addPost} from "./BoardsService.ts";
import {Toast} from "primereact/toast";
import {INewBoardPost} from "./INewBoardPost.ts";
import {NavigateFunction, useNavigate} from "react-router";

interface AddEditPostFormProps {
    post?: IBoardPost;
    boardSlug: string;
    parentId?: number;
    addPostCallback?: () => void;
}

export default function AddEditPostForm({post, boardSlug, parentId, addPostCallback}: AddEditPostFormProps) {
    let addPost$: Subject<IBoardPost>;
    const boardSlugRef = useRef<string>(boardSlug);
    const [isValid, setIsValid] = useState<boolean>(false);
    const toast: RefObject<Toast | null> = useRef<Toast>(null);
    const [postTitle, setPostTitle] = useState<string>("");
    const [postText, setPostText] = useState<string>("");
    const navigate: NavigateFunction = useNavigate();
    const navigateToNewPost = (newPost: IBoardPost) => {
        navigate(`/boards/${boardSlug}/posts/${newPost.slug}`);
    };
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
        const valid: boolean = validate();
        const post: INewBoardPost = {
            title: postTitle,
            postText: postText
        }
        if (parentId) {
            post.parentId = parentId;
        }
        if (valid) {
            addPost$ = addPost(post, boardSlugRef.current);
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
                    navigateToNewPost(newPost);
                    if (addPostCallback) {
                        addPostCallback();
                    }
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
        try {
            PostSchema.parse({
                title: postTitle,
                postText: postText,
            });
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

    return (
        <>
            <form onSubmit={onSubmit} className="w-full m-0 p-0">
                <div className="w-full mb-4">
                    <label className="mb-2 block" htmlFor="post-title">Title</label>
                    <InputText
                        className="w-full"
                        onChange={(e) => {
                            setPostTitle(e.target.value);
                            validate();
                        }}
                        value={postTitle}
                        maxLength={150}
                        invalid={!!formErrata.postTitle}
                        id="post-title"/>
                    <p className="text-red-500 pt-2">{formErrata.postTitle}</p>
                </div>

                <div className="w-full">
                    <label className="mb-2 block" htmlFor="post-text">Post text</label>
                    <InputTextarea
                        className="w-full"
                        onChange={(e) => {
                            setPostText(e.target.value);
                            validate();
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