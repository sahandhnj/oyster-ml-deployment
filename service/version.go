package service

import (
	"fmt"
	"path"
	"path/filepath"
	"strconv"

	"github.com/sahandhnj/apiclient/docker"

	"github.com/sahandhnj/apiclient/db"
	"github.com/sahandhnj/apiclient/filemanager"
	"github.com/sahandhnj/apiclient/types"
)

const (
	RequirementsFile     = "requirements.txt"
	ModelTarFile         = "model.tar.gz"
	DockerFileName       = "DockerFile"
	DockerIgnoreFileName = ".dockerignore"
	BuildLogFile         = "buildlog"
	RunFile              = "run.sh"
	TmpServerFile        = "server.py"
)

type VersionService struct {
	Version   *types.Version
	Model     *types.Model
	file      *filemanager.FileStoreManager
	DBHandler *db.DBStore
}

func NewVersionService(model *types.Model, dbHandler *db.DBStore) (*VersionService, error) {
	file, err := filemanager.NewFileStoreManager()
	if err != nil {
		return nil, err
	}

	versionSerimport os, io, sys, inspect, time, json, yaml, uuid
	import redis
	from PIL import Image
	from werkzeug.utils import secure_filename
	from flask import Flask, request, jsonify, flash
	import numpy as np
	# sys.path.insert(1, os.path.join(sys.path[0], '../..'))  # insert mlpipe root to path
	mlpipe_root = os.path.abspath("..")  
	sys.path.insert(0, mlpipe_root)
	
	
	
	# Set multiple paths for run testing period
	
	from config.clistyle import bcolor
	from servers.helperfunctions import base64_encoding, get_dtype
	
	try:
		from model import preprocessing as prepmod
	
		if hasattr(prepmod, 'preprocessing') and inspect.isfunction(prepmod.preprocessing):
			from model.preprocessing import preprocessing
			print("Preprocessing file available and loaded into vessel.")
		else:
			raise TypeError("Preprocessing file inserted, but does not contain function called 'preprocessing'.")
	except (ImportError):
		print("No preprocessing file inserted.")
	
	with open(mlpipe_root + "/settings.yaml", 'r') as stream:
		try:
			settings = yaml.load(stream)
		except yaml.YAMLError as exc:
			print(exc)
	
	
	
	# numpy.random.seed(42)
	app = Flask(__name__)
	rdb = redis.StrictRedis(
		host=settings['redis']['host'],
		port=settings['redis']['port'],
		db=settings['redis']['db']
	)
	
	# rdb.flushall()
	
	
	def allowed_file(filename):
		return '.' in filename and filename.rsplit('.', 1)[1].lower() in set(settings['data_stream']['allowed_extensions'])
	
	
	def get_file_type(filename):
		return '.' in filename and filename.rsplit('.', 1)[1].lower()
	
	
	@app.route('/predict', methods=["POST"])
	def predict():
	
		data = {"success": False}
	
		if request.method == "POST":
			# Check if file is inputted
			if 'data' not in request.files:
				flash("No file part")
				raise ValueError("No file part.")
			file = request.files['data']    ### Redundant?
			# print("FILENAME: ", file.filename)
			filetype = get_file_type(file.filename)
			# print("FILETYPE: ", filetype)
			# Check if file name is not empty
			if file.filename == '':
				flash("No selected file")
			if file and allowed_file(file.filename):
				filename = secure_filename(file.filename)
				if request.files.get('data'):
					user_input = request.files["data"].read()
					if (filetype in ['jpg', 'jpeg', 'png']):
						user_input = Image.open(io.BytesIO(user_input))
					else:
						pass
	
					preprocessed_input = preprocessing(user_input)
					# Get file properties
					if filetype in ['jpg', 'jpeg', 'png']:
						fileshape = np.array(preprocessed_input).shape
					else:
						fileshape = preprocessed_input.shape
	
					array_dtype = get_dtype(preprocessed_input)
					preprocessed_input = preprocessed_input.copy(order="C")
					encoded_input = base64_encoding(preprocessed_input)
	
					k = str(uuid.uuid4())
					d = {
						"id": k,
						"filename": filename,
						"filetype": filetype,
						"shape": fileshape,
						"dtype": array_dtype,
						"data": encoded_input
					}
					rdb.rpush(settings['data_stream']['data_queue'], json.dumps(d))  # dump the preprocessed input as a numpy array
	
					while True:
						output = rdb.get(k)
						if output is not None:
							output = output.decode("utf-8")
							# print("SUMMARY: ", json.loads(output)[0])
							data["summary"] = json.loads(output)
							rdb.delete(k)
							break
	
						time.sleep(settings['data_stream']['client_sleep'])
					data["success"] = True
	   
		return jsonify(data)
	
	
	@app.route("/predict")
	def hello():
		return "Hello, Welcome to Oysterbox Machine Learning Deployment!"
	
	
	if __name__ == "__main__":
		print((bcolor.BOLD + "* Loading Keras model and Flask starting server... \n"
			   "please wait until server has fully started" + bcolor.END))
		# print("* Starting model service... ")
		# t = Thread(target=classify_process, args=())
		# t.daemon
		# t.start()
		
		print("* Starting web service...")
		app.secret_key = settings['flask']['secret_key']
		# app.config['SESSION_TYPE'] = 'filesystem'
		app.run(
			host=settings['flask']['host'],
			port=int(settings['flask']['port']),
			debug=settings['flask']['debug']
		)vice := VersionService{
		file:      file,
		DBHandler: dbHandler,
		Model:     model,
	}

	return &versionService, nil
}

func (vs *VersionService) NewVersion() error {
	versionNumber := 0
	port := 5000

	versions, err := vs.DBHandler.VersionService.VersionsByModelId(vs.Model.ID)
	if err != nil {
		return err
	}

	for _, v := range versions {
		if v.VersionNumber > versionNumber {
			versionNumber = v.VersionNumber
		}
		if v.Port > port {
			port = v.Port
		}
	}
	versionNumber = versionNumber + 1
	port = port + 1

	version, err := types.NewVersion(versionNumber, vs.Model.ID)
	if err != nil {
		return err
	}
	version.ID = vs.DBHandler.VersionService.GetNextIdentifier()
	version.Port = port

	err = vs.DBHandler.VersionService.CreateVersion(version)
	if err != nil {
		return err
	}

	vs.Version = version

	err = vs.Apply()
	if err != nil {
		return err
	}

	return nil
}

func (vs *VersionService) PrintVersions() error {
	versions, err := vs.DBHandler.VersionService.VersionsByModelId(vs.Model.ID)
	if err != nil {
		return err
	}

	fmt.Printf("%s\t%s\t%s\t%s\n", "Name", "Version number", "Deployed", "Image Tag")
	for _, ver := range versions {
		fmt.Printf("%s\t%d\t%t\t\t%s\n", ver.Name, ver.VersionNumber, ver.Deployed, ver.ImageTag)
	}

	return nil
}

func (vs *VersionService) Deploy(versionNumber int, dcli *docker.DockerCli, verbose bool) error {
	version, err := vs.DBHandler.VersionService.VersionByVersionNumber(versionNumber, vs.Model.ID)
	if err != nil {
		return err
	}

	dockerFilePath := path.Join(vs.file.GetStorePath(version.Name), DockerFileName)
	mainTag := "oyster/" + vs.Model.Name + ":" + strconv.Itoa(version.VersionNumber)
	tags := []string{mainTag}

	fmt.Println("Deploying: ")
	fmt.Println(tags)

	logs, err := dcli.BuildImage(dockerFilePath, tags)
	if err != nil {
		return err
	}

	logFilePath := path.Join(vs.file.GetStorePath(version.Name), BuildLogFile)

	fmt.Println("Writing image build logs into: ", logFilePath)
	err = vs.file.WriteToFileWithReader(logFilePath, logs)
	if err != nil {
		return err
	}

	if verbose {
		err = vs.file.StreamFileToStdOut(logFilePath)
		if err != nil {
			return err
		}
	}

	version.ImageTag = mainTag
	containerName := vs.Model.Name + "" + strconv.Itoa(version.VersionNumber) + "-api"

	mountPath, err := filepath.Abs(vs.file.GetStorePath(version.Name))
	if err != nil {
		return err
	}

	containerId, err := dcli.CreateContainer(containerName, version.ImageTag, mountPath, strconv.Itoa(version.Port))
	if err != nil {
		return err
	}

	version.ContainerId = containerId

	if version.RedisEnabled {
		// redisContainerName := vs.Model.Name + "" + strconv.Itoa(version.VersionNumber) + "-redis"
		redisContainerName := "redis-generic"

		redisContainerId, err := dcli.CreateRedisContainer(redisContainerName)
		if err != nil {
			return err
		}

		version.RedisContainerId = redisContainerId

		networkName := vs.Model.Name + "" + strconv.Itoa(version.VersionNumber) + "-network"
		networkId, err := dcli.CreateNetwork(networkName)
		if err != nil {
			return err
		}

		version.NetworkId = networkId

		dcli.ConnectToNetwork(networkId, containerId)
		dcli.ConnectToNetwork(networkId, redisContainerId)
	}

	vs.DBHandler.VersionService.UpdateVersion(version.ID, version)

	return nil
}

func (vs *VersionService) Start(versionNumber int, dcli *docker.DockerCli) error {
	version, err := vs.DBHandler.VersionService.VersionByVersionNumber(versionNumber, vs.Model.ID)
	if err != nil {
		return err
	}

	if version.RedisEnabled {
		dcli.ContainerStart(version.RedisContainerId)
	}

	dcli.ContainerStart(version.ContainerId)

	return nil
}

func (vs *VersionService) Stop(versionNumber int, dcli *docker.DockerCli) error {
	version, err := vs.DBHandler.VersionService.VersionByVersionNumber(versionNumber, vs.Model.ID)
	if err != nil {
		return err
	}

	if version.RedisEnabled {
		dcli.ContainerStop(version.RedisContainerId)
	}

	dcli.ContainerStop(version.ContainerId)

	return nil
}

func (vs *VersionService) Down(versionNumber int, dcli *docker.DockerCli) error {
	version, err := vs.DBHandler.VersionService.VersionByVersionNumber(versionNumber, vs.Model.ID)
	if err != nil {
		return err
	}

	if version.RedisEnabled {
		dcli.ContainerDelete(version.RedisContainerId)
		version.RedisContainerId = ""
	}

	dcli.ContainerDelete(version.ContainerId)
	version.ContainerId = ""

	if version.NetworkId != "" {
		dcli.NetworkDelete(version.NetworkId)
		version.NetworkId = ""
	}

	vs.DBHandler.VersionService.UpdateVersion(version.ID, version)

	return nil
}

func (vs *VersionService) Apply() error {
	fm, err := filemanager.NewFileStoreManager()
	if err != nil {
		return err
	}

	fm.CreateDirectoryInStore(vs.Version.Name)
	fm.CTarGz(path.Join(vs.Version.Name, ModelTarFile), []string{vs.Model.ModelPath}, false)
	fm.CopyToStore(path.Join(vs.Model.ModelPath, RequirementsFile), path.Join(vs.Version.Name, RequirementsFile))
	fm.CopyToStore(path.Join(vs.Model.ModelPath, TmpServerFile), path.Join(vs.Version.Name, TmpServerFile))
	fm.CopyToStore(path.Join(vs.Model.ModelPath, RunFile), path.Join(vs.Version.Name, RunFile))

	vs.createDockerFile()

	return nil
}

func (vs *VersionService) createDockerFile() {
	docker_file_static = docker_file_static + "RUN pip3 install --user " + vs.file.ReadRQLineByLine(path.Join(vs.Version.Name, RequirementsFile))
	// docker_file_static = docker_file_static + "\nEXPOSE " + strconv.Itoa(vs.Version.Port) + "\n"
	docker_file_static = docker_file_static + "\nCMD bash run.sh"

	vs.file.WriteToFile(path.Join(vs.Version.Name, DockerFileName), docker_file_static)
	vs.file.WriteToFile(path.Join(vs.Version.Name, DockerIgnoreFileName), "")
}

var docker_file_static = `FROM python:3.6
ENV MODELPATH /src

RUN pip3 install --upgrade pip

WORKDIR $MODELPATH
`

var docker_file_static_big = `FROM ubuntu:18.04
ENV MODELPATH /src

RUN apt-get update && apt-get install -y software-properties-common

RUN add-apt-repository ppa:jonathonf/python-3.6
RUN apt-get update && apt-get install -y python3.6 curl python3-pip python-dev build-essential

RUN pip3 install --upgrade pip

WORKDIR $MODELPATH
`
