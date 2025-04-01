interface IProductReviewFormProps {
    productId: number;
}

export function ProductReviewForm() {

    const onSubmit = () => {

    };

    return (
        <>
            <form onSubmit={onSubmit}>
                <div className="flex flex-col gap-4">
                    <section className="">
                        <label htmlFor="title">Title</label>
                    </section>
                </div>
            </form>
        </>
    )
}