"""Initial schema

Revision ID: 001
Revises:
Create Date: 2025-04-12 02:46:00.000000

"""

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = "001"
down_revision = None
branch_labels = None
depends_on = None


def upgrade() -> None:
    # Create companies table
    op.create_table(
        "company",
        sa.Column("id", sa.Integer(), nullable=False),
        sa.Column("name", sa.String(length=255), nullable=False),
        sa.Column("careers_page_url", sa.String(length=255), nullable=False),
        sa.Column("logo_url", sa.String(length=255), nullable=True),
        sa.Column("description", sa.Text(), nullable=True),
        sa.Column("industry", sa.String(length=100), nullable=True),
        sa.Column(
            "active", sa.Boolean(), nullable=False, server_default=sa.text("true")
        ),
        sa.Column(
            "created_at", sa.DateTime(), nullable=False, server_default=sa.text("now()")
        ),
        sa.Column(
            "updated_at", sa.DateTime(), nullable=False, server_default=sa.text("now()")
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_company")),
    )
    op.create_index(op.f("idx_companies_name"), "company", ["name"], unique=False)
    op.create_index(
        op.f("idx_companies_industry"), "company", ["industry"], unique=False
    )

    # Create technologies table
    op.create_table(
        "technology",
        sa.Column("id", sa.Integer(), nullable=False),
        sa.Column("name", sa.String(length=100), nullable=False),
        sa.Column("category", sa.String(length=50), nullable=True),
        sa.Column("parent_id", sa.Integer(), nullable=True),
        sa.Column(
            "created_at", sa.DateTime(), nullable=False, server_default=sa.text("now()")
        ),
        sa.ForeignKeyConstraint(
            ["parent_id"],
            ["technology.id"],
            name=op.f("fk_technology_parent_id_technology"),
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_technology")),
        sa.UniqueConstraint("name", name=op.f("uq_technology_name")),
    )

    # Create jobs table
    op.create_table(
        "job",
        sa.Column("id", sa.Integer(), nullable=False),
        sa.Column("company_id", sa.Integer(), nullable=True),
        sa.Column("title", sa.String(length=255), nullable=False),
        sa.Column("slug", sa.String(length=255), nullable=False),
        sa.Column("description", sa.Text(), nullable=False),
        sa.Column("requirements", sa.Text(), nullable=False),
        sa.Column("preferred_skills", sa.Text(), nullable=True),
        sa.Column("experience_level", sa.String(length=50), nullable=False),
        sa.Column("employment_type", sa.String(length=50), nullable=False),
        sa.Column("location", sa.String(length=50), nullable=True),
        sa.Column("work_mode", sa.String(length=20), nullable=False),
        sa.Column("application_url", sa.String(length=255), nullable=True),
        sa.Column("job_function", sa.String(length=100), nullable=True),
        sa.Column(
            "first_seen_at",
            sa.DateTime(),
            nullable=False,
            server_default=sa.text("now()"),
        ),
        sa.Column(
            "last_seen_at",
            sa.DateTime(),
            nullable=False,
            server_default=sa.text("now()"),
        ),
        sa.Column("posted_at", sa.DateTime(), nullable=True),
        sa.Column(
            "is_active", sa.Boolean(), nullable=False, server_default=sa.text("true")
        ),
        sa.Column("signature", sa.String(length=64), nullable=False),
        sa.Column(
            "created_at", sa.DateTime(), nullable=False, server_default=sa.text("now()")
        ),
        sa.Column(
            "updated_at", sa.DateTime(), nullable=False, server_default=sa.text("now()")
        ),
        sa.ForeignKeyConstraint(
            ["company_id"], ["company.id"], name=op.f("fk_job_company_id_company")
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_job")),
        sa.UniqueConstraint("signature", name=op.f("uq_job_signature")),
    )

    # Create indexes for jobs table
    op.create_index(
        op.f("idx_jobs_title_tsvector"),
        "job",
        [sa.text("to_tsvector('english', title)")],
        unique=False,
        postgresql_using="gin",
    )
    op.create_index(
        op.f("idx_jobs_description_tsvector"),
        "job",
        [sa.text("to_tsvector('english', description)")],
        unique=False,
        postgresql_using="gin",
    )
    op.create_index(op.f("idx_jobs_location"), "job", ["location"], unique=False)
    op.create_index(
        op.f("idx_active_jobs"),
        "job",
        ["id"],
        unique=False,
        postgresql_where=sa.text("is_active = true"),
    )
    op.create_index(op.f("idx_jobs_work_mode"), "job", ["work_mode"], unique=False)
    op.create_index(op.f("idx_jobs_posted_at"), "job", ["posted_at"], unique=False)
    op.create_index(op.f("idx_jobs_company_id"), "job", ["company_id"], unique=False)
    op.create_index(
        op.f("idx_jobs_job_function"), "job", ["job_function"], unique=False
    )
    op.create_index(
        op.f("idx_jobs_employment_type"), "job", ["employment_type"], unique=False
    )
    op.create_index(
        op.f("idx_jobs_experience_level"), "job", ["experience_level"], unique=False
    )

    # Create job_technologies junction table
    op.create_table(
        "job_technology",
        sa.Column("id", sa.Integer(), nullable=False),
        sa.Column("job_id", sa.Integer(), nullable=True),
        sa.Column("technology_id", sa.Integer(), nullable=True),
        sa.Column(
            "is_primary", sa.Boolean(), nullable=False, server_default=sa.text("false")
        ),
        sa.Column(
            "created_at", sa.DateTime(), nullable=False, server_default=sa.text("now()")
        ),
        sa.ForeignKeyConstraint(
            ["job_id"], ["job.id"], name=op.f("fk_job_technology_job_id_job")
        ),
        sa.ForeignKeyConstraint(
            ["technology_id"],
            ["technology.id"],
            name=op.f("fk_job_technology_technology_id_technology"),
        ),
        sa.PrimaryKeyConstraint("id", name=op.f("pk_job_technology")),
    )

    # Create indexes for job_technologies table
    op.create_index(
        op.f("idx_job_technologies_job_id"), "job_technology", ["job_id"], unique=False
    )
    op.create_index(
        op.f("idx_job_technologies_technology_id"),
        "job_technology",
        ["technology_id"],
        unique=False,
    )


def downgrade() -> None:
    # Drop tables in reverse order of creation
    op.drop_table("job_technology")
    op.drop_table("job")
    op.drop_table("technology")
    op.drop_table("company")
