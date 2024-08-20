-- 1. Roles
CREATE TABLE roles (
    ID SERIAL PRIMARY KEY,
    role_name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 2. Permissions
CREATE TABLE permissions (
    ID SERIAL PRIMARY KEY,
    permission_desc VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 3. Users
CREATE TABLE users (
    ID SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    pin VARCHAR(255),
    merchant_id VARCHAR(255),
    role_id INT REFERENCES roles(ID),
    user_type VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 4. Role Permissions
CREATE TABLE role_permissions (
    ID SERIAL PRIMARY KEY,
    role_id INT REFERENCES roles(ID),
    permission_id INT REFERENCES permissions(ID),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 5. Payment Methods
CREATE TABLE payment_methods (
    ID SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    pay_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 6. Merchants
CREATE TABLE merchants (
    ID SERIAL PRIMARY KEY,
    merchant_id VARCHAR(255) UNIQUE NOT NULL,
    merchant_name VARCHAR(255) UNIQUE NOT NULL,
    merchant_secret VARCHAR(255) UNIQUE NOT NULL,
    currency VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 7. Merchant Accounts
CREATE TABLE merchant_accounts (
    ID SERIAL PRIMARY KEY,
    merchant_id VARCHAR(255) REFERENCES merchants(merchant_id),
    settle_balance DECIMAL(18,2) DEFAULT 0.00,
    not_settle_balance DECIMAL(18,2) DEFAULT 0.00,
    hold_balance DECIMAL(18,2) DEFAULT 0.00,
    balance_capital_flow DECIMAL(18,2) DEFAULT 0.00,
    pending_transaction_out DECIMAL(18,2) DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 8. Merchant Payment Methods
CREATE TABLE merchant_payment_methods (
    ID SERIAL PRIMARY KEY,
    merchant_id INT REFERENCES merchants(ID),
    payment_method_id INT REFERENCES payment_methods(ID),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 9. Providers
CREATE TABLE providers (
    ID SERIAL PRIMARY KEY,
    provider_id VARCHAR(255) UNIQUE NOT NULL,
    provider_name VARCHAR(255) NOT NULL,
    currency VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 10. Provider Payment Methods
CREATE TABLE provider_payment_methods (
    ID SERIAL PRIMARY KEY,
    provider_id INT REFERENCES providers(ID),
    payment_method_id INT REFERENCES payment_methods(ID),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 11. Bank Lists
CREATE TABLE bank_lists (
    ID SERIAL PRIMARY KEY,
    bank_name VARCHAR(255) NOT NULL,
    bank_code VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 12. Merchant Payment Method Bank Lists
CREATE TABLE merchant_payment_method_bank_lists (
    ID SERIAL PRIMARY KEY,
    merchant_payment_method_id INT REFERENCES merchant_payment_methods(ID),
    bank_list_id INT REFERENCES bank_lists(ID),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 13. Provider Payment Method Bank Lists
CREATE TABLE provider_payment_method_bank_lists (
    ID SERIAL PRIMARY KEY,
    provider_payment_method_id INT REFERENCES provider_payment_methods(ID),
    bank_list_id INT REFERENCES bank_lists(ID),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 14. Merchant Paychannels
CREATE TABLE merchant_paychannels (
    ID SERIAL PRIMARY KEY,
    merchant_payment_method_id INT REFERENCES merchant_payment_methods(ID),
    segment VARCHAR(255),
    fee DECIMAL(18,2) NOT NULL DEFAULT 0.00,
    fee_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    min_transaction DECIMAL(18,2) DEFAULT 0.00,
    max_transaction DECIMAL(18,2) DEFAULT 0.00,
    max_daily_transaction DECIMAL(18,2) DEFAULT 0.00,
    merchant_paychannel_code VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 15. Provider Paychannels
CREATE TABLE provider_paychannels (
    ID SERIAL PRIMARY KEY,
    provider_payment_method_id INT REFERENCES provider_payment_methods(ID),
    paychannel_name VARCHAR(255),
    fee DECIMAL(18,2) NOT NULL DEFAULT 0.00,
    fee_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    min_transaction DECIMAL(18,2) DEFAULT 0.00,
    max_transaction DECIMAL(18,2) DEFAULT 0.00,
    max_daily_transaction DECIMAL(18,2) DEFAULT 0.00,
    interface_setting VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 16. Provider Paychannel Bank List
CREATE TABLE provider_paychannel_bank_lists (
    ID SERIAL PRIMARY KEY,
    provider_paychannel_id INT REFERENCES provider_paychannels(ID),
    bank_list_id INT REFERENCES bank_lists(ID),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 17. Transactions
CREATE TABLE transactions (
    ID SERIAL PRIMARY KEY,
    payment_id VARCHAR(255) UNIQUE NOT NULL,
    merchant_reference_number VARCHAR(255) UNIQUE NOT NULL,
    provider_reference_number VARCHAR(255) UNIQUE,
    merchant_paychannel_id INT REFERENCES merchant_paychannels(ID),
    provider_paychannel_id INT REFERENCES provider_paychannels(ID),
    transaction_amount DECIMAL(18,2) NOT NULL DEFAULT 0.00,
    bank_code VARCHAR(50),
    status VARCHAR(50) NOT NULL,
    client_ip_address VARCHAR(255),
    merchant_callback_url VARCHAR(255),
    request_method VARCHAR(255) NOT NULL,
    transaction_payment_generated TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 18. Merchant Callbacks
CREATE TABLE merchant_callbacks (
    ID SERIAL PRIMARY KEY,
    payment_id VARCHAR(255) REFERENCES transactions(payment_id),
    callback_status VARCHAR(50) NOT NULL,
    payment_status_in_callback VARCHAR(50) NOT NULL,
    callback_result TEXT,
    triggered_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 19. Account Informations
CREATE TABLE account_informations (
    ID SERIAL PRIMARY KEY,
    payment_id VARCHAR(255) NOT NULL REFERENCES transactions(payment_id),
    account_number VARCHAR(255),
    account_name VARCHAR(255),
    bank_name VARCHAR(255),
    bank_code VARCHAR(50),
    reference_number VARCHAR(255),
    remark TEXT,
    phone_number VARCHAR(50),
    email VARCHAR(255),
    account_type VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 20. Provider Transaction Confirmation Details
CREATE TABLE provider_transaction_confirmation_details (
    ID SERIAL PRIMARY KEY,
    payment_id VARCHAR(255) REFERENCES transactions(payment_id),
    type VARCHAR(50) NOT NULL,
    confirmation_result TEXT NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 21. Transaction Status Logs
CREATE TABLE transaction_status_logs (
    ID SERIAL PRIMARY KEY,
    payment_id VARCHAR(255) REFERENCES transactions(payment_id),
    status_log TEXT NOT NULL,
    change_by VARCHAR(255) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 22. Paychannel Routings
CREATE TABLE paychannel_routings (
    ID SERIAL PRIMARY KEY,
    provider_paychannel_id INT REFERENCES provider_paychannels(ID),
    merchant_paychannel_id INT REFERENCES merchant_paychannels(ID),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 23. Routing Histories
CREATE TABLE histories_operations (
    ID SERIAL PRIMARY KEY,
    history_type VARCHAR(50) NOT NULL,
    payload_before VARCHAR(255),
    payload_after VARCHAR(255),
    username VARCHAR(50),
    activity VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 24. Merchant Paychannel Daily Transactions
CREATE TABLE merchant_paychannel_daily_transactions (
    ID SERIAL PRIMARY KEY,
    merchant_paychannel_id INT REFERENCES merchant_paychannels(ID),
    transaction_amount DECIMAL(18,2) DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 25. Provider Paychannel Daily Transactions
CREATE TABLE provider_paychannel_daily_transactions (
    ID SERIAL PRIMARY KEY,
    provider_paychannel_id INT REFERENCES provider_paychannels(ID),
    transaction_amount DECIMAL(18,2) DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 26. Merchant Capital Flows
CREATE TABLE merchant_capital_flows (
    ID SERIAL PRIMARY KEY,
    payment_id VARCHAR(255) NOT NULL,
    merchant_account_id INT NOT NULL,
    temp_balance DECIMAL(18,2) NOT NULL DEFAULT 0.00,
    transaction_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    capital_type VARCHAR(50),
    notes VARCHAR(255),
    reverse_from VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 27. Provider Credentials
CREATE TABLE provider_credentials (
    ID SERIAL PRIMARY KEY,
    provider_id VARCHAR(255) REFERENCES providers(provider_id),
    key VARCHAR(255),
    value VARCHAR(255),
    interface_setting VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 28. Report Storages
CREATE TABLE report_storages (
    ID SERIAL PRIMARY KEY,
    merchant_id VARCHAR(50),
    period VARCHAR(255),
    export_type VARCHAR(50),
    status VARCHAR(50),
    report_url VARCHAR(255),
    created_by_user VARCHAR(50),
    file_name VARCHAR(255) UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
