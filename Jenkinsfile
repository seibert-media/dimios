def label = "buildpod.${env.JOB_NAME}".replaceAll(/[^A-Za-z-]+/, '-').take(62) + "p"

podTemplate(
	name: label,
	label: label,
	containers: [
		containerTemplate(
			name: 'build-golang',
			image: 'docker.io/bborbe/build-golang:1.0.0',
			ttyEnabled: true,
			command: 'cat',
			resourceRequestCpu: '500m',
			resourceRequestMemory: '500Mi',
			resourceLimitCpu: '2000m',
			resourceLimitMemory: '500Mi',
		),
		containerTemplate(
			name: 'build-docker',
			image: 'docker.io/bborbe/build-docker:1.0.0',
			ttyEnabled: true,
			command: 'cat',
			privileged: true,
			resourceRequestCpu: '500m',
			resourceRequestMemory: '500Mi',
			resourceLimitCpu: '2000m',
			resourceLimitMemory: '500Mi',
		),
	],
	volumes: [
		secretVolume(mountPath: '/home/jenkins/.docker', secretName: 'docker'),
		hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock'),
	],
	inheritFrom: '',
	namespace: 'jenkins',
	serviceAccount: '',
	workspaceVolume: emptyDirWorkspaceVolume(false),
) {
	node(label) {
		properties([
			buildDiscarder(logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '3', numToKeepStr: '5')),
			pipelineTriggers([
				cron('H 2 * * *'),
				pollSCM('H/5 * * * *'),
			]),
			parameters([
				string(name: 'Version', defaultValue: '', description: 'Version to build'),
			]),
		])
		try {
			container('build-golang') {
				stage('Golang Checkout') {
					timeout(time: 5, unit: 'MINUTES') {
						checkout([
							$class: 'GitSCM',
							branches: scm.branches,
							doGenerateSubmoduleConfigurations: scm.doGenerateSubmoduleConfigurations,
							extensions: scm.extensions + [[$class: 'CloneOption', noTags: false, reference: '', shallow: true]],
							submoduleCfg: [],
							userRemoteConfigs: scm.userRemoteConfigs
						])
					}
				}
				stage('Golang Link') {
					timeout(time: 5, unit: 'MINUTES') {
						sh """
						mkdir -p /go/src/github.com/bborbe
						ln -s `pwd` /go/src/github.com/bborbe/k8s_deploy
						"""
					}
				}
				stage('Golang Test') {
					timeout(time: 15, unit: 'MINUTES') {
						sh "cd /go/src/github.com/bborbe/k8s_deploy && make test"
					}
				}
			}
			container('build-docker') {
				stage('Docker Checkout') {
					timeout(time: 5, unit: 'MINUTES') {
						checkout([
							$class: 'GitSCM',
							branches: scm.branches,
							doGenerateSubmoduleConfigurations: scm.doGenerateSubmoduleConfigurations,
							extensions: scm.extensions + [[$class: 'CloneOption', noTags: false, reference: '', shallow: true]],
							submoduleCfg: [],
							userRemoteConfigs: scm.userRemoteConfigs
						])
					}
				}
				stage('Docker Trigger') {
					timeout(time: 5, unit: 'MINUTES') {
						env.TRIGGER = sh (script: "VERSION=${params.Version} make trigger", returnStdout: true).trim()
						echo "trigger = ${env.TRIGGER}"
					}
				}
				stage('Docker Build') {
					if (env.TRIGGER == 'build' || env.BRANCH_NAME != 'master') {
						timeout(time: 15, unit: 'MINUTES') {
							sh "VERSION=${params.Version} make build"
						}
					}
				}
				stage('Docker Upload') {
					if (env.TRIGGER == 'build' && env.BRANCH_NAME == 'master') {
						timeout(time: 15, unit: 'MINUTES') {
							sh "VERSION=${params.Version} make upload"
						}
					}
				}
				stage('Docker Clean') {
					if (env.TRIGGER == 'build' || env.BRANCH_NAME != 'master') {
						timeout(time: 5, unit: 'MINUTES') {
							sh "VERSION=${params.Version} make clean"
						}
					}
				}
			}
			currentBuild.result = 'SUCCESS'
		} catch (any) {
			currentBuild.result = 'FAILURE'
			throw any //rethrow exception to prevent the build from proceeding
		} finally {
			if ('FAILURE'.equals(currentBuild.result)) {
				emailext(
					body: '${DEFAULT_CONTENT}',
					mimeType: 'text/html',
					replyTo: '$DEFAULT_REPLYTO',
					subject: '${DEFAULT_SUBJECT}',
					to: emailextrecipients([
						[$class: 'CulpritsRecipientProvider'],
						[$class: 'RequesterRecipientProvider']
					])
				)
			}
		}
	}
}
