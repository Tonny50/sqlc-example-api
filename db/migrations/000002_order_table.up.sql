CREATE TABLE "order" (
    "customer_id" VARCHAR(36) NOT NULL,
    "product_name" VARCHAR(30) NOT NULL,
    "price" VARCHAR(10) NOT NULL,
    "tran_status" VARCHAR(20) NOT NULL,
    "order_date" TIMESTAMP DEFAULT now() NOT NULL,
    FOREIGN KEY ("customer_id")
    REFERENCES "customer"("id")
    ON DELETE CASCADE

    
)