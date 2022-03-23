# TODO -> Balance
# - Amount
# - Currency [ Currency Table ]
# - Owner [User table]
# - Status [int 1: active ,2: pending]
# - Moderation bool
# - CreatedAt
# - UpdatedAt

# TODO -> Reaction
# - PostType (UID)
# - Reaction Type [Int 1, 2, 3, 4, 5, 6, 7]
# - Owner [User table]
# - CreatedAt
# - UpdatedAt

# TODO -> Address Table
# - Product [Product Table]
# - Name
# - Address 1
# - Address 2
# - City
# - State
# - Country
# - Latitude
# - Longitude
# - Log IP
# - Status [int 1,2,3]
# - Moderation [int 1,2,3]
# - User [User table]
# - Business [User table]
# - Post [ Post table ]
# - Comment [ Comment table ]
# - Order [ Order table ]
# - Review [ Review table ]
# - Qna [ Qna table ]

# TODO -> Category Fields
# - Name
# - Icon
# - Moderation [int 1,2,3]
# - Status [int 1,2,3]
# - Product [Product Table]
# - Author [User table] <admin | end_user>
# - Business [User Table]
# - Child [Category table] <subcategory>

# TODO -> Group Fields
# - Name
# - Icon
# - About
# - Status [int 1,2,3]
# - Moderation [int 1,2,3]
# - Owner [User]
# - Editor [User]
# - Publisher [User]
# - Child [ Group table ] <subgroup>
# - CreatedAt
# - UpdatedAt

# TODO -> QnA
# - Title
# - Description
# - Tags [Hashtag table]
# - Status [int 1,2,3]
# - Moderation [int 1,2,3]
# - Privacy [int 1,2,3]
# - Author [User]
# - Address [Address table]
# - CreatedAt
# - UpdatedAt
# - Comments [Comments table]
# - Reaction [Reaction table]

# TODO -> Product
# - Title
# - Regular Price
# - Selling Price
# - Currency [ Currency table]
# - Color
# - Size
# .....
# - Hashtags [Hashtag table]
# - Category [Category Table]
# - Variants [Product Table]
# - Order [Order table]
# - Review [Review table]
# - Type [int 1,2,3] ~ physical, digital, services
# - Supported [int 1: supported,2: not supported,3: null]
# - Downloadable [Assets table]
# - Thumbnail [Assets table]
# - Gallery [ Assets table ]
# - Excerpt
# - Description
# - Technical Information
# - Additional Information
# - Product Information
# - Product Guides
# - Status [int 1,2,3]
# - Moderation [int 1,2,3]
# - Address [Address Table]
# - Owner [Business table] required
# - Author [User table]
# - Category [Category Table]

# TODO -> Review
# - Business [User Table]
# - Product [Product Table]
# - Rating [int 1,2,3,4,5]
# - Title
# - Description
# - Status [int 1,2]
# - Author [User Table]
# - Assets [Assets Table]
# - Moderation [int 1,2,3]

# TODO -> Assets (Image/Video)
# - Path
# - Type ( e.g, 1,2.. ( 1=> Image, 2=>Video))
# - CreatedAt
# - UpdatedAt
# - Owner [User table]
# - Business [ Business table]
# - Product [ Product table]
# - Post [ Post table ]
# - Comment [ Comment table ]
# - Review [ Review table ]
# - Qna [ Qna Table]
# - Moderation bool

# Only businesses can sell products
# TODO ORDER table
# - Post [ Post table ]
# - Product [ Product table ]
# - TransactionID
# - Sender [User table]
# - Business [User table]
# - Receiver [User table]
# - Status [int 1,2,3]
# - Moderation bool
# - Coupon [Coupon table]
# - Amount [float]
# - Currency [Currency table]
# - CreatedAt
# - UpdatedAt
# - UserAgent