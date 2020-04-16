package gradesapi

// This file just contains types for Canvas responses generated by
// ChimeraCoder/gojson.

type canvasEnrollmentType string

const (
	canvasEnrollmentTypeStudentEnrollment  = canvasEnrollmentType("student")
	canvasEnrollmentTypeObserverEnrollment = canvasEnrollmentType("observer")
)

// /api/v1/users/:userID
type canvasUserProfileResponse canvasUserProfile

type canvasUserProfile struct {
	AvatarURL string      `json:"avatar_url"`
	Bio       interface{} `json:"bio"`
	Calendar  struct {
		Ics string `json:"ics"`
	} `json:"calendar"`
	EffectiveLocale string      `json:"effective_locale"`
	ID              uint64      `json:"id"`
	IntegrationID   interface{} `json:"integration_id"`
	Locale          interface{} `json:"locale"`
	LoginID         string      `json:"login_id"`
	LtiUserID       string      `json:"lti_user_id"`
	Name            string      `json:"name"`
	PrimaryEmail    string      `json:"primary_email"`
	ShortName       string      `json:"short_name"`
	SortableName    string      `json:"sortable_name"`
	TimeZone        string      `json:"time_zone"`
	Title           interface{} `json:"title"`
}

// /api/v1/users/:userID/observees
type canvasUserObserveesResponse []canvasObservee

// canvasObservee is an Observee in Canvas
// no link because they don't have a published example struct
type canvasObservee struct {
	CreatedAt                     string  `json:"created_at"`
	ID                            uint64  `json:"id"`
	Name                          string  `json:"name"`
	ObservationLinkRootAccountIds []int64 `json:"observation_link_root_account_ids"`
	ShortName                     string  `json:"short_name"`
	SortableName                  string  `json:"sortable_name"`
}

// /api/v1/courses
type canvasCoursesResponse []canvasCourse

// canvasCourse represents a Course in Canvas.
// https://canvas.instructure.com/doc/api/courses.html#Course
type canvasCourse struct {
	AccountID                   int64 `json:"account_id"`
	ApplyAssignmentGroupWeights bool  `json:"apply_assignment_group_weights"`
	Blueprint                   bool  `json:"blueprint"`
	Calendar                    struct {
		Ics string `json:"ics"`
	} `json:"calendar"`
	CourseCode                       string             `json:"course_code"`
	CreatedAt                        string             `json:"created_at"`
	DefaultView                      string             `json:"default_view"`
	EndAt                            string             `json:"end_at"`
	EnrollmentTermID                 int64              `json:"enrollment_term_id"`
	Enrollments                      []canvasEnrollment `json:"enrollments"`
	GradePassbackSetting             interface{}        `json:"grade_passback_setting"`
	GradingStandardID                int64              `json:"grading_standard_id"`
	HideFinalGrades                  bool               `json:"hide_final_grades"`
	ID                               uint64             `json:"id"`
	IsPublic                         bool               `json:"is_public"`
	IsPublicToAuthUsers              bool               `json:"is_public_to_auth_users"`
	License                          string             `json:"license"`
	Locale                           string             `json:"locale"`
	Name                             string             `json:"name"`
	OverriddenCourseVisibility       string             `json:"overridden_course_visibility"`
	PublicSyllabus                   bool               `json:"public_syllabus"`
	PublicSyllabusToAuth             bool               `json:"public_syllabus_to_auth"`
	RestrictEnrollmentsToCourseDates bool               `json:"restrict_enrollments_to_course_dates"`
	RootAccountID                    int64              `json:"root_account_id"`
	StartAt                          string             `json:"start_at"`
	StorageQuotaMb                   int64              `json:"storage_quota_mb"`
	TimeZone                         string             `json:"time_zone"`
	UUID                             string             `json:"uuid"`
	WorkflowState                    string             `json:"workflow_state"`

	CanvasCBLHidden bool `json:"canvascbl_hidden"`
}

// canvasEnrollment represents a user's enrollment in a course.
// https://canvas.instructure.com/doc/api/enrollments.html#Enrollment
type canvasEnrollment struct {
	AssociatedUserID               uint64               `json:"associated_user_id"`
	EnrollmentState                string               `json:"enrollment_state"`
	LimitPrivilegesToCourseSection bool                 `json:"limit_privileges_to_course_section"`
	Role                           string               `json:"role"`
	RoleID                         int64                `json:"role_id"`
	Type                           canvasEnrollmentType `json:"type"`
	UserID                         uint64               `json:"user_id"`
	ComputedCurrentScore           float64              `json:"computed_current_score"`
	ComputedCurrentGrade           string               `json:"computed_current_grade"`
}

// /api/v1/courses/:courseID/outcome_rollups
type canvasOutcomeRollupsResponse struct {
	Meta struct {
		Pagination struct {
			Count     int64  `json:"count"`
			Current   string `json:"current"`
			First     string `json:"first"`
			Last      string `json:"last"`
			Page      int64  `json:"page"`
			PageCount int64  `json:"page_count"`
			PerPage   int64  `json:"per_page"`
			Template  string `json:"template"`
		} `json:"pagination"`
	} `json:"meta"`
	Rollups []canvasOutcomeRollup `json:"rollups"`
}

// canvasOutcomeRollup is an Outcome rollup in Canvas
// https://canvas.instructure.com/doc/api/outcome_results.html#OutcomeResult
type canvasOutcomeRollup struct {
	Links struct {
		Section string `json:"section"`
		User    string `json:"user"`
	} `json:"links"`
	Scores []struct {
		Count      int64 `json:"count"`
		HidePoints bool  `json:"hide_points"`
		Links      struct {
			Outcome string `json:"outcome"`
		} `json:"links"`
		Score       float64 `json:"score"`
		SubmittedAt string  `json:"submitted_at"`
		Title       string  `json:"title"`
	} `json:"scores"`
}

// /api/v1/courses/:courseID/assignments
type canvasAssignmentsResponse []canvasAssignment

// canvasAssignment is an assignment in Canvas
// https://canvas.instructure.com/doc/api/assignments.html#Assignment
type canvasAssignment struct {
	AllowedAttempts                int64  `json:"allowed_attempts"`
	AnonymizeStudents              bool   `json:"anonymize_students"`
	AnonymousGrading               bool   `json:"anonymous_grading"`
	AnonymousInstructorAnnotations bool   `json:"anonymous_instructor_annotations"`
	AnonymousPeerReviews           bool   `json:"anonymous_peer_reviews"`
	AnonymousSubmissions           bool   `json:"anonymous_submissions"`
	AssignmentGroupID              int64  `json:"assignment_group_id"`
	AutomaticPeerReviews           bool   `json:"automatic_peer_reviews"`
	CanDuplicate                   bool   `json:"can_duplicate"`
	CourseID                       int64  `json:"course_id"`
	CreatedAt                      string `json:"created_at"`
	//Description                     string      `json:"description"`
	DueAt                           string      `json:"due_at"`
	DueDateRequired                 bool        `json:"due_date_required"`
	FinalGraderID                   interface{} `json:"final_grader_id"`
	FreeFormCriterionComments       bool        `json:"free_form_criterion_comments"`
	GradeGroupStudentsIndividually  bool        `json:"grade_group_students_individually"`
	GraderCommentsVisibleToGraders  bool        `json:"grader_comments_visible_to_graders"`
	GraderCount                     int64       `json:"grader_count"`
	GraderNamesVisibleToFinalGrader bool        `json:"grader_names_visible_to_final_grader"`
	GradersAnonymousToGraders       bool        `json:"graders_anonymous_to_graders"`
	GradingStandardID               interface{} `json:"grading_standard_id"`
	GradingType                     string      `json:"grading_type"`
	GroupCategoryID                 interface{} `json:"group_category_id"`
	HasSubmittedSubmissions         bool        `json:"has_submitted_submissions"`
	HTMLURL                         string      `json:"html_url"`
	ID                              uint64      `json:"id"`
	InClosedGradingPeriod           bool        `json:"in_closed_grading_period"`
	IntraGroupPeerReviews           bool        `json:"intra_group_peer_reviews"`
	IsQuizAssignment                bool        `json:"is_quiz_assignment"`
	LockAt                          string      `json:"lock_at"`
	LockExplanation                 string      `json:"lock_explanation"`
	LockInfo                        struct {
		AssetString string `json:"asset_string"`
		CanView     bool   `json:"can_view"`
		LockAt      string `json:"lock_at"`
	} `json:"lock_info"`
	LockedForUser          bool        `json:"locked_for_user"`
	MaxNameLength          int64       `json:"max_name_length"`
	ModeratedGrading       bool        `json:"moderated_grading"`
	Muted                  bool        `json:"muted"`
	Name                   string      `json:"name"`
	OmitFromFinalGrade     bool        `json:"omit_from_final_grade"`
	OnlyVisibleToOverrides bool        `json:"only_visible_to_overrides"`
	OriginalAssignmentID   interface{} `json:"original_assignment_id"`
	OriginalAssignmentName interface{} `json:"original_assignment_name"`
	OriginalCourseID       interface{} `json:"original_course_id"`
	OriginalQuizID         interface{} `json:"original_quiz_id"`
	PeerReviews            bool        `json:"peer_reviews"`
	PointsPossible         float64     `json:"points_possible"`
	Position               int64       `json:"position"`
	PostManually           bool        `json:"post_manually"`
	PostToSis              bool        `json:"post_to_sis"`
	Published              bool        `json:"published"`
	QuizID                 int64       `json:"quiz_id"`
	Rubric                 []struct {
		CriterionUseRange bool `json:"criterion_use_range"`
		//Description       string  `json:"description"`
		ID               string `json:"id"`
		IgnoreForScoring bool   `json:"ignore_for_scoring"`
		//LongDescription   string  `json:"long_description"`
		OutcomeID int64   `json:"outcome_id"`
		Points    float64 `json:"points"`
		Ratings   []struct {
			//Description     string  `json:"description"`
			ID string `json:"id"`
			//LongDescription string  `json:"long_description"`
			Points float64 `json:"points"`
		} `json:"ratings"`
		VendorGUID interface{} `json:"vendor_guid"`
	} `json:"rubric"`
	RubricSettings struct {
		FreeFormCriterionComments bool    `json:"free_form_criterion_comments"`
		HidePoints                bool    `json:"hide_points"`
		HideScoreTotal            bool    `json:"hide_score_total"`
		ID                        int64   `json:"id"`
		PointsPossible            float64 `json:"points_possible"`
		Title                     string  `json:"title"`
	} `json:"rubric_settings"`
	SecureParams           string   `json:"secure_params"`
	SubmissionTypes        []string `json:"submission_types"`
	SubmissionsDownloadURL string   `json:"submissions_download_url"`
	UnlockAt               string   `json:"unlock_at"`
	UpdatedAt              string   `json:"updated_at"`
	UseRubricForGrading    bool     `json:"use_rubric_for_grading"`
	WorkflowState          string   `json:"workflow_state"`
}

// /api/v1/courses/:courseID/outcome_results
type canvasOutcomeResultsResponse struct {
	OutcomeResults []canvasOutcomeResult `json:"outcome_results"`
}

type canvasOutcomeResult struct {
	Hidden     bool   `json:"hidden"`
	HidePoints bool   `json:"hide_points"`
	ID         uint64 `json:"id"`
	Links      struct {
		Alignment       string `json:"alignment"`
		Assignment      string `json:"assignment"`
		LearningOutcome string `json:"learning_outcome"`
		User            string `json:"user"`
	} `json:"links"`
	Mastery               bool    `json:"mastery"`
	Percent               float64 `json:"percent"`
	Possible              float64 `json:"possible"`
	Score                 float64 `json:"score"`
	SubmittedOrAssessedAt string  `json:"submitted_or_assessed_at"`
}

// /api/v1/outcomes/:outcomeID
type canvasOutcomeResponse canvasOutcome

// CanvasOutcome is an Outcome in Canvas
// https://canvas.instructure.com/doc/api/outcomes.html#Outcome
type canvasOutcome struct {
	Assessed          bool   `json:"assessed"`
	CalculationInt    int64  `json:"calculation_int"`
	CalculationMethod string `json:"calculation_method"`
	CanEdit           bool   `json:"can_edit"`
	ContextID         uint64 `json:"context_id"`
	ContextType       string `json:"context_type"`
	//Description          string  `json:"description"`
	DisplayName          string  `json:"display_name"`
	HasUpdateableRubrics bool    `json:"has_updateable_rubrics"`
	ID                   uint64  `json:"id"`
	MasteryPoints        float64 `json:"mastery_points"`
	PointsPossible       float64 `json:"points_possible"`
	Ratings              []struct {
		//Description string  `json:"description"`
		Points float64 `json:"points"`
	} `json:"ratings"`
	Title      string      `json:"title"`
	URL        string      `json:"url"`
	VendorGUID interface{} `json:"vendor_guid"`
}

type canvasOutcomeAlignmentsResponse []canvasOutcomeAlignment

type canvasOutcomeAlignment struct {
	AssignmentID      int64  `json:"assignment_id"`
	LearningOutcomeID int64  `json:"learning_outcome_id"`
	SubmissionTypes   string `json:"submission_types"`
	Title             string `json:"title"`
	URL               string `json:"url"`
}

// POST /login/oauth2/token with ?grant_type=refresh_token
type canvasRefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
	User        struct {
		EffectiveLocale string `json:"effective_locale"`
		GlobalID        string `json:"global_id"`
		ID              int64  `json:"id"`
		Name            string `json:"name"`
	} `json:"user"`
}

type canvasErrorArrayResponse struct {
	Errors []struct {
		Message string `json:"message"`
	}
}

type canvasErrorResponse struct {
	Error string `json:"error"`
}

func (cer canvasErrorResponse) toCanvasErrorArrayResponse() canvasErrorArrayResponse {
	return canvasErrorArrayResponse{
		Errors: []struct {
			Message string `json:"message"`
		}{
			{Message: cer.Error},
		},
	}
}

type canvasOAuth2ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
