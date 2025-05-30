-- 000006_add_brief_system.down.sql
-- Migration Down: Remove brief system
-- Safely removes all brief-related tables, indexes, and types in reverse dependency order
-- Preserves existing schema (users, companies, contacts, datasources, sales_processes)

-- Step 1: Drop Performance Indexes
-- Remove all indexes created for the brief system to avoid dependency issues
DROP INDEX IF EXISTS idx_sales_process_briefs_sales_process_id; -- Sales process integration index
DROP INDEX IF EXISTS idx_ground_truth_field_name; -- Ground truth field lookup index
DROP INDEX IF EXISTS idx_ground_truth_master_brief_id; -- Ground truth aggregation index
DROP INDEX IF EXISTS idx_sales_process_tracking_brief_id; -- Sales tracking (fields 141-150) index
DROP INDEX IF EXISTS idx_behavioral_insights_brief_id; -- Behavioral insights (fields 126-140) index
DROP INDEX IF EXISTS idx_technical_integration_brief_id; -- Technical integration (fields 111-125) index
DROP INDEX IF EXISTS idx_project_requirements_brief_id; -- Project requirements (fields 96-110) index
DROP INDEX IF EXISTS idx_financial_procurement_brief_id; -- Financial procurement (fields 81-95) index
DROP INDEX IF EXISTS idx_competitive_intelligence_brief_id; -- Competitive intelligence (fields 66-80) index
DROP INDEX IF EXISTS idx_current_state_assessment_brief_id; -- Current state (fields 51-65) index
DROP INDEX IF EXISTS idx_buying_committee_brief_id; -- Buying committee (fields 31-50) index
DROP INDEX IF EXISTS idx_strategic_context_brief_id; -- Strategic context (fields 16-30) index
DROP INDEX IF EXISTS idx_company_intelligence_brief_id; -- Company intelligence (fields 1-15) index
DROP INDEX IF EXISTS idx_brief_attachments_brief_id; -- Attachment lookup index
DROP INDEX IF EXISTS idx_briefs_master_brief_id; -- Brief hierarchy index
DROP INDEX IF EXISTS idx_master_briefs_contact_id; -- Master brief contact lookup index
DROP INDEX IF EXISTS idx_master_briefs_company_id; -- Master brief company lookup index
DROP INDEX IF EXISTS idx_master_briefs_cognito_sub; -- Master brief user ownership index

-- Step 2: Drop Junction and Integration Tables
-- Remove tables that link brief system to existing sales processes
DROP TABLE IF EXISTS sales_process_briefs; -- Links briefs to sales_processes table

-- Step 3: Drop Ground Truth System
-- Remove aggregated intelligence data storage
DROP TABLE IF EXISTS ground_truth; -- Master aggregated values with confidence scoring

-- Step 4: Drop Intelligence Tables (in reverse dependency order)
-- Remove all 9 intelligence category tables containing the 150 fields
DROP TABLE IF EXISTS sales_process_tracking; -- Fields 141-150: Sales funnel and deal progression
DROP TABLE IF EXISTS behavioral_insights; -- Fields 126-140: Decision patterns and communication styles
DROP TABLE IF EXISTS technical_integration; -- Fields 111-125: IT architecture and integration requirements
DROP TABLE IF EXISTS project_requirements; -- Fields 96-110: Implementation and project management
DROP TABLE IF EXISTS financial_procurement; -- Fields 81-95: Budget and procurement processes
DROP TABLE IF EXISTS competitive_intelligence; -- Fields 66-80: Competitor analysis and evaluation
DROP TABLE IF EXISTS current_state_assessment; -- Fields 51-65: Existing solutions and pain points
DROP TABLE IF EXISTS buying_committee; -- Fields 31-50: Decision makers and influencers
DROP TABLE IF EXISTS strategic_context; -- Fields 16-30: Corporate strategy and initiatives
DROP TABLE IF EXISTS company_intelligence; -- Fields 1-15: Basic company information and market position

-- Step 5: Drop Core Brief System Tables
-- Remove brief content and attachment management
DROP TABLE IF EXISTS brief_attachments; -- Media attachments linked to datasources
DROP TABLE IF EXISTS briefs; -- Individual briefs with content and metadata
DROP TABLE IF EXISTS master_briefs; -- Top-level brief containers

-- Step 6: Drop Custom Enum Types
-- Remove brief-specific enumeration types
DROP TYPE IF EXISTS brief_tag; -- Processing stage tags (initial, specific, updated, final)
DROP TYPE IF EXISTS brief_type; -- Brief categories (master, stage_specific, regular)