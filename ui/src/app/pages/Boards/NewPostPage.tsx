import AddEditPostForm from "../../components/Boards/AddEditPostForm.tsx";

export default function NewPostPage() {
    return (
        <>
            <h1 className="text-3xl font-bold mb-4">New Post</h1>
            <section className="mt-4 w-1/2">
                <AddEditPostForm/>
            </section>
        </>
    )
}