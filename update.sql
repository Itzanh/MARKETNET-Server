ALTER TABLE public."user"
    ADD COLUMN uses_google_authenticator boolean NOT NULL DEFAULT false;

ALTER TABLE public."user"
    ADD COLUMN google_authenticator_secret character(8);